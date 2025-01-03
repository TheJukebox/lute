package lute_api

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	streamPb "lute/gen/stream"
	uploadPb "lute/gen/upload"
	"lute/internal/convert"
	apiErrors "lute/internal/lute_api/errors"

	"github.com/google/uuid"
)

// placeholder while we work out locking and such
var uploads = make(map[string]*uploadPb.UploadRequest)

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

// This handler is for validating that the file CAN be uploaded
// so we want to do things like create locks here!
func (s *UploadService) StartUpload(_ context.Context, request *uploadPb.UploadRequest) (*uploadPb.UploadResponse, error) {
	// generate an ID
	file_id := uuid.NewString()
	uploads[file_id] = request
	log.Printf("(%v) Received a request to begin upload!", file_id)

	os.Create(fmt.Sprintf("uploads/raw/%v", request.GetFileName()))

	return &uploadPb.UploadResponse{
		FileId: file_id,
	}, nil
}

func (s *UploadService) UploadChunk(_ context.Context, chunk *uploadPb.Chunk) (*uploadPb.ChunkResponse, error) {
	// Check an upload request was made first
	file_id := chunk.GetFileId()
	request, ok := uploads[file_id]
	if !ok {
		log.Printf("No upload request found for %v", chunk.GetFileId())
		return nil, &apiErrors.NoUploadRequest{
			FileId: file_id,
		}
	}

	filename := fmt.Sprintf("uploads/raw/%v", request.GetFileName())
	output_path := strings.Split(request.GetFileName(), ".")[0] + ".aac"
	output_path = fmt.Sprintf("uploads/converted/%v", output_path)
	output, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		log.Printf("(%v) failed to upload file: '%v':\n\t%v", file_id, filename, err)
		return nil, err
	}
	_, err = output.Write(chunk.GetData())
	if err != nil {
		log.Printf("(%v) Failed to write chunk: %v", file_id, err)
		return nil, err
	}
	log.Printf("(%v) Got chunk of size: %v", file_id, len(chunk.GetData()))
	if chunk.GetFinal() {
		log.Printf("(%v) Final chunk received for file: '%v'", file_id, filename)
		converted_path, err := convert.ConvertFile(filename, output_path)
		if err != nil {
			log.Printf("(%v) Upload failed: %v", file_id, err)
		}
		log.Printf("(%v) Finished conversion: %v", file_id, converted_path)
	}
	return &uploadPb.ChunkResponse{
		Success: true,
		Message: "Received chunk!",
	}, nil
}

// func (s *UploadService) FileUpload(_ context.Context, file *uploadPb.FileUploadRequest) (*uploadPb.FileUploadResponse, error) {
// 	log.Printf("Received a file upload request...")
// 	filename := file.GetFileName()
// 	data := file.GetFileData()
//
// 	// Validate the filename
// 	validName, err := forbiddenChars(filename)
// 	if err != nil {
// 		log.Printf("Unable to parse filename: '%v'", filename)
// 		return nil, err
// 	} else if !validName {
// 		log.Printf("Invalid filename: '%v'", filename)
// 		return nil, &apiErrors.IllegalFileName{Filename: filename}
// 	}
//
// 	// We should ensure that the data is the right format too.
// 	// We'll have to investigate how best to do that - magic bytes?
//
// 	// Try to write the file
// 	filename = fmt.Sprintf("uploads/raw/%v", filename)
// 	outputPath := fmt.Sprintf("uploads/converted/%v", strings.Split(filename, "/")[2])
// 	output, err := os.Create(filename)
// 	if err != nil {
// 		log.Printf("Could not open file for write: %s\n", filename)
// 		return &uploadPb.FileUploadResponse{Success: false, Message: "File failed to upload: could not open file for write"}, err
// 	}
// 	defer output.Close()
//
// 	_, err = output.Write(data)
// 	if err != nil {
// 		log.Printf("Could not write file: %s\n", filename)
// 		return &uploadPb.FileUploadResponse{Success: false, Message: "Failed to write file."}, err
// 	}
//
// 	// Perform file conversion
// 	result, err := convert.ConvertFile(filename, outputPath)
// 	if err != nil {
// 		log.Printf("Failed to convert file '%v'! Cleaning up...", filename)
// 		os.Remove(filename)
// 		return nil, err
// 	}
// 	log.Printf("Created file '%v'", result)
//
// 	return &uploadPb.FileUploadResponse{Success: true, Message: "Successfully uploaded"}, nil
// }
//
