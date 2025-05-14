package middleware

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"io"
	"log"
	"net/http"

	streamPb "lute/gen/stream"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

func CreateGrpcClient() (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	conn.Connect()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")

			if r.Method == "OPTIONS" {
				log.Printf("(%s) Handling CORS request...", r.Header.Get("Origin"))
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		},
	)
}

func GrpcWebParseMiddleware(grpcServer *grpc.Server, next http.Handler, client streamPb.AudioStreamClient) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
            // handling the wire format
            // https://protobuf.dev/programming-guides/encoding/
            if r.Header.Get("Content-Type") == "application/lute-grpc" {
                origin := r.Header.Get("Origin")
                log.Printf("(%s) Parsing incoming stream request...", origin)
                
                // read the body
                body, err := io.ReadAll(r.Body)
                if err != nil {
                    log.Printf("Failed to read request: %v", err)
                    http.Error(w, "Failed to read request", http.StatusBadRequest)
                    return
                }

                // log the body
                log.Printf("body: %v", body)

                // unmarshal the data
                var msg streamPb.AudioStreamRequest
                err = proto.Unmarshal(body, &msg)
                if err != nil {
                    log.Printf("Failed to unmarshal request")
                }
                filename := msg.FileName
                session_id := msg.SessionId
				log.Printf("(%s) (%s) Requesting audio stream: %s", origin, session_id, filename)

                stream, err := client.StreamAudio(context.Background(), &msg)
                for {
                    data , err := stream.Recv()
                    if err == io.EOF {
                        break
                    }
					if err != nil {
						log.Printf("(%s) Couldn't finish stream: %q", origin, err)
						http.Error(w, "Error sending stream chunk", http.StatusInternalServerError)
						return
					}

                    // create the chunk
                    chunk := &streamPb.AudioStreamChunk{
                        Data: data.GetData(),
                        Sequence: data.GetSequence(),
                    }
                    // encode the chunk into wire format
                    encodedResponse, err := frameGrpcResponse(chunk)
					if err != nil {
						log.Printf("(%s) Failed to frame data: %q", origin, err)
						http.Error(w, "Failed to frame response", http.StatusInternalServerError)
						return
					}
                    
                    // TODO: make this output debug only
                    log.Printf("(%s) Streaming chunk (%d bytes)...", origin, len(encodedResponse))
                    w.Write(encodedResponse)

                    // push the chunk to the user
                    if flusher, ok := w.(http.Flusher); ok {
                        flusher.Flush()
                    }
                }

                // close stream
				log.Printf("(%s) Sending gRPC trailers to client...", origin)
				grpc.SendHeader(r.Context(), metadata.New(map[string]string{
					"grpc-status":  "0",
					"grpc-message": "OK",
				}))
				log.Printf("(%s) Stream complete!", origin)
				return
            }
			// Only operate on requests from grpc-web clients
			if r.Header.Get("Content-Type") == "application/grpc-web-text" {
				origin := r.Header.Get("Origin")
				log.Printf("(%s) Parsing incoming grpc-web-text request...", origin)
				// Read body
				body, err := io.ReadAll(r.Body)
				if err != nil {
					log.Printf("Failed to read request: %v", err)
					http.Error(w, "Failed to read request", http.StatusBadRequest)
					return
				}

				// Decode from b64
				decodedBody := make([]byte, base64.StdEncoding.DecodedLen(len(body)))
				n, err := base64.StdEncoding.Decode(decodedBody, body)
				if err != nil {
					log.Printf("Failed to decode request: %v", err)
					http.Error(w, "Failed to decode request", http.StatusInternalServerError)
					return
				}
				decodedBody = decodedBody[5:n] // trim the frame

				// Unmarshal into protobuf
				var msg streamPb.AudioStreamRequest
				if err := proto.Unmarshal(decodedBody, &msg); err != nil {
					log.Printf("Failed to unmarshal protobuf: %v", err)
					http.Error(w, "Failed to unmarshall protobuf", http.StatusInternalServerError)
					return
				}
				filename := msg.FileName
				session_id := msg.SessionId
				log.Printf("(%s) (%s) Requesting audio stream: %s", origin, session_id, filename)

				// this is slowing us down and making response time about 10s in total...
				// not sure how to speed it up, we're kind of just making 2 requests, so we have to wait for the 2nd one
				// to resolve before we can respond to the first.
				stream, err := client.StreamAudio(context.Background(), &msg)
				if err != nil {
					log.Printf("Couldn't start stream: %v", err)
					http.Error(w, "Failed to start stream", http.StatusInternalServerError)
					return
				}
				log.Printf("(%s) Opened gRPC stream...", origin)

				w.Header().Set("Content-Type", "application/grpc-web-text")
				w.Header().Set("Cache-Control", "no-cache")
				w.Header().Set("Connection", "keep-alive")

				log.Printf("(%s) Headers set, starting stream...", origin)
				for {
					// Receive the next part of the stream from the gRPC server
					data, err := stream.Recv()
					if err == io.EOF {
						break
					}
					if err != nil {
						log.Printf("(%s) Couldn't finish stream: %q", origin, err)
						http.Error(w, "Error sending stream chunk", http.StatusInternalServerError)
						return
					}

					// Serialize the chunk for sending to the client
					chunk := &streamPb.AudioStreamChunk{
						Data:     data.GetData(),
						Sequence: data.GetSequence(),
					}
					encodedResponse, err := frameGrpcResponse(chunk)
					if err != nil {
						log.Printf("(%s) Failed to frame data: %q", origin, err)
						http.Error(w, "Failed to frame response", http.StatusInternalServerError)
						return
					}

					// write the chunk data
					w.Write(encodedResponse)
					if flusher, ok := w.(http.Flusher); ok {
						flusher.Flush()
					}
				}

				// Terminate the stream
				log.Printf("(%s) Sending gRPC trailers to client...", origin)
				grpc.SendHeader(r.Context(), metadata.New(map[string]string{
					"grpc-status":  "0",
					"grpc-message": "OK",
				}))
				log.Printf("(%s) Stream complete!", origin)
				return
			}
		},
	)
}

func frameGrpcResponse(data *streamPb.AudioStreamChunk) ([]byte, error) {
	// Marshal data into protobuf format
	dataBytes, err := proto.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Create a wire compliant frame
	frameLength := uint32(len(dataBytes)) // we have to add the length
	var frameLengthBuffer bytes.Buffer
	if err := binary.Write(&frameLengthBuffer, binary.BigEndian, frameLength); err != nil {
		return nil, err
	}
    frame := []byte{0x00}                               // compression flag
	frame = append(frame, frameLengthBuffer.Bytes()...) // the length in big endian

	// frame the data
	// 0x00 | 0x00 0x00 0x00 0x00 | DATA
	response := append(frame, dataBytes...)
	// base64 encode the response
	//encodedResponse := make([]byte, base64.StdEncoding.EncodedLen(len(response)))
	//base64.StdEncoding.Encode(encodedResponse, response)

	return response, nil
}
