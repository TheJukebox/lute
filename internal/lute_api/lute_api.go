package lute_api

import (
	"context"
	"io"
	"log"
	"os"

	streamPb "lute/gen/stream"
	uploadPb "lute/gen/upload"
)

type StreamService struct {
	streamPb.UnimplementedAudioStreamServer
}

type UploadService struct {
	uploadPb.UnimplementedUploadServer
}

func (s *StreamService) StreamAudio(request *streamPb.AudioStreamRequest, stream streamPb.AudioStream_StreamAudioServer) error {
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

func (s *UploadService) FileUpload(_ context.Context, in *uploadPb.FileUploadRequest) (*uploadPb.FileUploadResponse, error) {
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
