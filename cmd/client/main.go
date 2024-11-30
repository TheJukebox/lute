package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"

	pb "lute/gen/stream" // Replace with your generated package import path
)

func main() {
	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewAudioStreamClient(conn)

	// Request to stream a specific audio file
	req := &pb.AudioStreamRequest{
		FileName:  "output.aac",
		SessionId: "session-123",
	}

	// Stream audio from the server
	stream, err := client.StreamAudio(context.Background(), req)
	if err != nil {
		log.Fatalf("Error starting audio stream: %v", err)
	}

	// Process the received chunks
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			log.Println("Audio stream ended")
			break
		}
		if err != nil {
			log.Fatalf("Error receiving audio chunk: %v", err)
		}

		// Log the received chunk
		fmt.Printf("Received chunk %d with %d bytes\n", chunk.GetSequence(), len(chunk.GetData()))
	}
}
