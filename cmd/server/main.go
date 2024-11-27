package main

import (
	"context"
	"log"
	"net"

	pb "lute/gen/converter"

	"google.golang.org/grpc"
)

type converter struct {
	pb.UnimplementedConverterServiceServer
}

func (s *converter) ConvertToHLS(_ context.Context, in *pb.FileUploadRequest) (*pb.FileUploadResponse, error) {
	log.Printf("Received File Upload Request: %v", in.GetFileName())
	return &pb.FileUploadResponse{Success: true, Message: "Successfully uploaded"}, nil
}

func main() {
	log.Println("Starting Lute server...")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterConverterServiceServer(s, &converter{})
	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
