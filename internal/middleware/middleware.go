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

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")

			if r.Method == "OPTIONS" {
				log.Printf("Received CORS request from %s", r.Header.Get("Origin"))
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		},
	)
}

func GrpcWebParseMiddleware(grpcServer *grpc.Server, next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Only operate on requests from grpc-web clients
			if r.Header.Get("Content-Type") == "application/grpc-web-text" {
				log.Printf("Received grpc-web-text request from %s", r.Header.Get("Origin"))

				// Read body
				body, err := io.ReadAll(r.Body)
				if err != nil {
					log.Printf("Failed to read request: %v", err)
					http.Error(w, "Failed to read request", http.StatusBadRequest)
					return
				}
				log.Printf("body: %s", string(body))

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
				log.Printf("Session %s (%s) is requesting stream: %s", session_id, r.Header.Get("Origin"), filename)

				conn, err := grpc.NewClient(
					"localhost:50051",
					grpc.WithTransportCredentials(insecure.NewCredentials()),
				)
				if err != nil {
					log.Printf("Failed to connect to gRPC server: %v", err)
					http.Error(w, "Failed to forward request", http.StatusInternalServerError)
					return
				}
				defer conn.Close()

				client := streamPb.NewAudioStreamClient(conn)

				log.Println("Connected to gRPC server.")
				stream, err := client.StreamAudio(context.Background(), &msg)
				if err != nil {
					log.Printf("Couldn't start stream: %v", err)
					http.Error(w, "Failed to start stream", http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/grpc-web-text")
				w.Header().Set("Cache-Control", "no-cache")
				w.Header().Set("Connection", "keep-alive")

				for {
					// Receive the next part of the stream from the gRPC server
					data, err := stream.Recv()
					if err == io.EOF {
						break
					}
					if err != nil {
						log.Printf("Couldn't finish stream: %v", err)
						http.Error(w, "Failed to finish stream", http.StatusInternalServerError)
						return
					}

					// Serialize the chunk for sending to the client
					chunk := &streamPb.AudioStreamChunk{
						Data:     data.GetData(),
						Sequence: data.GetSequence(),
					}
					chunkBytes, err := proto.Marshal(chunk)
					if err != nil {
						log.Printf("Unable to marshal chunk: %v", err)
						http.Error(w, "Failed to finish stream", http.StatusInternalServerError)
						return
					}
					var encodedChunkBytes = make([]byte, base64.StdEncoding.EncodedLen(len(chunkBytes)))
					base64.StdEncoding.Encode(encodedChunkBytes, chunkBytes)

					frameLength := uint32(len(encodedChunkBytes))
					var frameLengthBuffer bytes.Buffer
					if err := binary.Write(&frameLengthBuffer, binary.LittleEndian, frameLength); err != nil {
						log.Printf("Failed to create frame: %v", err)
						http.Error(w, "Failed while creating frame", http.StatusInternalServerError)
						return
					}
					frame := []byte{0x00}
					frame = append(frame, frameLengthBuffer.Bytes()...)
					response := append(frame, encodedChunkBytes...)
					log.Printf("FRAME: %05x | CALC_LENGTH: %d", frame[:5], frameLength)

					// write the chunk data
					w.Write(response)
					if flusher, ok := w.(http.Flusher); ok {
						flusher.Flush()
					}
				}
				// Terminate the stream
				grpc.SendHeader(r.Context(), metadata.New(map[string]string{
					"grpc-status":  "0",
					"grpc-message": "OK",
				}))
				log.Println("Finished stream.")
				return
			}
		},
	)
}
