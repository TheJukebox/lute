package main

import (
	"log"
	"net"
	"net/http"
	"os"

	streamPb "lute/gen/stream"
	uploadPb "lute/gen/upload"
	api "lute/internal/lute_api"
	mw "lute/internal/middleware"

	"google.golang.org/grpc"
)

func debugSetup() {
	log.Printf("Creating debug folders...")
	os.Mkdir("uploads/raw", 0700)
	os.Mkdir("uploads/converted", 0700)
}

func main() {
	// setup folders
	debugSetup()

	listener, _ := net.Listen("tcp", "127.0.0.1:50051")
	grpcNative := grpc.NewServer()
	uploadPb.RegisterUploadServer(grpcNative, &api.UploadService{})
	streamPb.RegisterAudioStreamServer(grpcNative, &api.StreamService{})
	go grpcNative.Serve(listener)
	log.Printf("gRPC Native server listening at %v", listener.Addr())

	grpcWeb := grpc.NewServer()
	uploadPb.RegisterUploadServer(grpcWeb, &api.UploadService{})
	streamPb.RegisterAudioStreamServer(grpcWeb, &api.StreamService{})

	// standup client for HTTP/1.1 to HTTP/2
	client, err := mw.CreateGrpcClient()
	if err != nil {
		log.Fatalf("Failed to standup internal gRPC client...")
	}
	defer client.Close()
	audioClient := streamPb.NewAudioStreamClient(client)
	log.Println("Started audio streaming client for HTTP/1.1 to HTTP/2...")

	mux := http.NewServeMux()
	mux.Handle("/upload.Upload/UploadFile", grpcWeb)
	mux.Handle("/stream.AudioStream/StreamAudio", grpcWeb)
	middleware := mw.GrpcWebParseMiddleware(grpcWeb, mux, audioClient)
	middleware = mw.CorsMiddleware(middleware)

	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: middleware,
	}

	log.Printf("Server listening at %v", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Failed to serve: %v", err)
	}
}
