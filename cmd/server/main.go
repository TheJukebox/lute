package main

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	streamPb "lute/gen/stream"
	uploadPb "lute/gen/upload"

	mw "lute/internal/middleware"

	"google.golang.org/grpc"
)

type streamService struct {
	streamPb.UnimplementedAudioStreamServer
}

type uploadService struct {
	uploadPb.UnimplementedUploadServer
}

func (s *streamService) StreamAudio(request *streamPb.AudioStreamRequest, stream streamPb.AudioStream_StreamAudioServer) error {
	file, err := os.Open(request.GetFileName())
	if err != nil {
		log.Printf("Failed to open file for streaming: %q\n", err)
		return err
	}
	defer file.Close()

	streamBuffer := make([]byte, 1024)
	sequence := int32(0)

	for {
		chunkSize, err := file.Read(streamBuffer)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("Failed to chunk stream: %q\n", err)
		}
		chunk := &streamPb.AudioStreamChunk{
			Data:     streamBuffer[:chunkSize],
			Sequence: sequence,
		}
		sequence++
		stream.Send(chunk)
	}
}

func (s *uploadService) FileUpload(_ context.Context, in *uploadPb.FileUploadRequest) (*uploadPb.FileUploadResponse, error) {
	log.Printf("Received File Upload Request: %v", in.GetFileName())
	output, err := os.Create(in.GetFileName())
	if err != nil {
		log.Printf("Could not open file for write: %s\n", in.GetFileName())
		return &uploadPb.FileUploadResponse{Success: false, Message: "File failed to upload: could not open file for write"}, err
	}
	defer output.Close()

	_, err = output.Write(in.GetFileData())
	if err != nil {
		log.Printf("Could not write file: %s\n", in.GetFileName())
		return &uploadPb.FileUploadResponse{Success: false, Message: "Failed to write file."}, err
	}
	return &uploadPb.FileUploadResponse{Success: true, Message: "Successfully uploaded"}, nil
}

func main() {
	listener, _ := net.Listen("tcp", "127.0.0.1:50051")
	grpcNative := grpc.NewServer()
	uploadPb.RegisterUploadServer(grpcNative, &uploadService{})
	streamPb.RegisterAudioStreamServer(grpcNative, &streamService{})
	go grpcNative.Serve(listener)
	log.Printf("gRPC Native server listening at %v", listener.Addr())

	grpcWeb := grpc.NewServer()
	uploadPb.RegisterUploadServer(grpcWeb, &uploadService{})
	streamPb.RegisterAudioStreamServer(grpcWeb, &streamService{})

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
