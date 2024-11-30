package main

import (
	"context"
	"log"
	"net"
	"os"

	pb "lute/gen/upload"

	"google.golang.org/grpc"
)

type uploadService struct {
	pb.UnimplementedUploadServer
}

func (s *uploadService) FileUpload(_ context.Context, in *pb.FileUploadRequest) (*pb.FileUploadResponse, error) {
	log.Printf("Received File Upload Request: %v", in.GetFileName())
	output, err := os.Create(in.GetFileName())
	if err != nil {
		log.Printf("Could not open file for write: %s\n", in.GetFileName())
		return &pb.FileUploadResponse{Success: false, Message: "File failed to upload: could not open file for write"}, err
	}
	defer output.Close()

	_, err = output.Write(in.GetFileData())
	if err != nil {
		log.Printf("Could not write file: %s\n", in.GetFileName())
		return &pb.FileUploadResponse{Success: false, Message: "Failed to write file."}, err
	}
	return &pb.FileUploadResponse{Success: true, Message: "Successfully uploaded"}, nil
}

func main() {
	log.Println("Starting Lute server...")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUploadServer(s, &uploadService{})
	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
