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

<<<<<<< HEAD
func debugSetup() {
	os.Mkdir("uploads/raw", 0700)
	os.Mkdir("uploads/converted", 0700)
=======
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

	streamBuffer := make([]byte, 8192*5)
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
>>>>>>> c9bc42aeb0a21c9618460309c831caf7669ab5f3
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
