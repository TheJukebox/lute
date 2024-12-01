package middleware

import (
	"bytes"
	"encoding/base64"
	"io"
	"log"
	"net/http"

	streamPb "lute/gen/stream"

	"google.golang.org/grpc"
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

				forwardMsg, err := proto.Marshal(&msg)
				if err != nil {
					log.Printf("Failed to marshal message: %v", err)
					http.Error(w, "Failed to marshall message", http.StatusInternalServerError)
					return
				}
				log.Printf("forward: %s", &forwardMsg)
				reader := bytes.NewReader(forwardMsg)

				forwardReq, err := http.NewRequest("POST", "stream.AudioStream/StreamAudio", reader)
				if err != nil {
					log.Printf("Failed to forward request: %v", err)
					http.Error(w, "Failed to forward request", http.StatusInternalServerError)
					return
				}
				forwardReq.Header.Set("Content-Type", "application/grpc")
				forwardReq.Header.Del("Content-Length")
				log.Println(forwardReq.Header)
				grpcServer.ServeHTTP(w, forwardReq)
				return
			}
			next.ServeHTTP(w, r)
		},
	)
}
