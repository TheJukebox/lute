package main

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	streamPb "lute/gen/stream"
	uploadPb "lute/gen/upload"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
)

type streamService struct {
	streamPb.UnimplementedAudioStreamServer
}

type uploadService struct {
	uploadPb.UnimplementedUploadServer
}

func (s *streamService) StreamAudio(request *streamPb.AudioStreamRequest, stream streamPb.AudioStream_StreamAudioServer) error {
	log.Printf("Received Stream Request from %s for %s", request.SessionId, request.FileName)
	file, err := os.Open(request.GetFileName())
	if err != nil {
		log.Printf("Something has gone terribly wrong: %q\n", err)
		return err
	}
	defer file.Close()

	streamBuffer := make([]byte, 1024)
	sequence := int32(0)

	for {
		chunkSize, err := file.Read(streamBuffer)
		if err == io.EOF {
			log.Println("Finished streaming!")
			return nil
		}
		if err != nil {
			log.Printf("Failed to chunk stream: %s\n", err)
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

func startGrpcServer(offline chan bool) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to bind port: %v", err)
	}
	s := grpc.NewServer()
	uploadPb.RegisterUploadServer(s, &uploadService{})
	streamPb.RegisterAudioStreamServer(s, &streamService{})

	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Printf("Failed to serve: %v", err)
		offline <- true
	}
	time.Sleep(0)
}

func startHttpServer(offline chan bool) {
	conn, err := grpc.NewClient("http://localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to start HTTP server: %v", err)
		offline <- true
	}
	defer conn.Close()

	mux := runtime.NewServeMux()
	streamPb.RegisterAudioStreamHandlerServer(context.Background(), mux, &streamService{})
	handler := cors.Default().Handler(mux)
	log.Println("Starting HTTP server at localhost:8080")
	http.ListenAndServe(":8080", handler)
}

func main() {
	offline := make(chan bool)
	log.Println("Starting Lute server...")
	go startGrpcServer(offline)
	go startHttpServer(offline)
	<-offline
}
