package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"

	streamPb "lute/gen/stream"
	uploadPb "lute/gen/upload"

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

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, x-grpc-web, x-user-agent")

			if r.Method == "OPTIONS" {
				log.Println("Handling a CORS request...")
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		},
	)
}

func main() {
	s := grpc.NewServer()
	uploadPb.RegisterUploadServer(s, &uploadService{})
	streamPb.RegisterAudioStreamServer(s, &streamService{})

	mux := http.NewServeMux()
	mux.Handle("/stream.AudioStream/StreamAudio", s)
	handlerWithCors := corsMiddleware(mux)

	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: handlerWithCors,
	}

	log.Printf("Server listening at %v", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Failed to serve: %v", err)
	}
}
