package middleware

import (
	"io"
	"log"
	"net/http"

	"google.golang.org/grpc"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, x-grpc-web, x-user-agent")

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
			if r.Header.Get("Content-Type") == "application/grpc-web-text" {
				body, err := io.ReadAll(r.Body)
				log.Printf("Received grpc-web-text request from %s", r.Header.Get("Origin"))
				if err != nil {
					http.Error(w, "Failed to read the request body", http.StatusBadRequest)
					return
				}
				log.Printf("Body:\n%s", body)
			}
			// grpcServer.ServeHTTP()
			next.ServeHTTP(w, r)
		},
	)
}
