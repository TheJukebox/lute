package lute_api

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"

	streamPb "lute/gen/stream"
	uploadPb "lute/gen/upload"
	"lute/internal/convert"
	apiErrors "lute/internal/lute_api/errors"
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

func forbiddenChars(s string) (bool, error) {
	r, err := regexp.Compile("^[a-zA-Z0-9._-]+$")
	if err != nil {
		return false, err
	}
	match := r.MatchString(s)
	return match, nil
}

func (s *UploadService) FileUpload(_ context.Context, file *uploadPb.FileUploadRequest) (*uploadPb.FileUploadResponse, error) {
	log.Printf("Received a file upload request...")
	filename := file.GetFileName()
	data := file.GetFileData()

	// Validate the filename
	validName, err := forbiddenChars(filename)
	if err != nil {
		log.Printf("Unable to parse filename: '%v'", filename)
		return nil, err
	} else if !validName {
		log.Printf("Invalid filename: '%v'", filename)
		return nil, &apiErrors.IllegalFileName{Filename: filename}
	}

	// We should ensure that the data is the right format too.
	// We'll have to investigate how best to do that - magic bytes?

	// Try to write the file
	filename = fmt.Sprintf("uploads/raw/%v", filename)
	outputPath := fmt.Sprintf("uploads/converted/%v", filename)
	output, err := os.Create(filename)
	if err != nil {
		log.Printf("Could not open file for write: %s\n", filename)
		return &uploadPb.FileUploadResponse{Success: false, Message: "File failed to upload: could not open file for write"}, err
	}
	defer output.Close()

	_, err = output.Write(data)
	if err != nil {
		log.Printf("Could not write file: %s\n", filename)
		return &uploadPb.FileUploadResponse{Success: false, Message: "Failed to write file."}, err
	}

	// Perform file conversion
	result, err := convert.ConvertFile(filename, outputPath)
	if err != nil {
		log.Printf("Failed to convert file '%v'! Cleaning up...", filename)
		os.Remove(filename)
		return nil, err
	}
	log.Printf("Created file '%v'", result)

	return &uploadPb.FileUploadResponse{Success: true, Message: "Successfully uploaded"}, nil
}
