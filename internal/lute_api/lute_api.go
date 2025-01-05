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
	// Check the validity of the filename
	valid, err := forbiddenChars(request.GetFileName())
	if err != nil {
		log.Printf("Failed to check the validity of the filename: '%v'", request.GetFileName())
		return nil, err
	}
	if !valid {
		log.Printf("Upload has invalid filename: '%v'", request.GetFileName())
		return nil, &apiErrors.IllegalFileName{
			Filename: request.GetFileName(),
		}
	}

	// generate an ID
	file_id := uuid.NewString()
	uploads[file_id] = request
	log.Printf("(%v) Received a request to begin upload!", file_id)

	os.Create(fmt.Sprintf("uploads/raw/%v", request.GetFileName()))

	// Return a file ID, ready for chunks to be written
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

	// create the paths for the file
	filename := fmt.Sprintf("uploads/raw/%v", request.GetFileName())
	output_path := strings.Split(request.GetFileName(), ".")[0] + ".aac"
	output_path = fmt.Sprintf("uploads/converted/%v", output_path)

	// Open the file for writing
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

	// Handle closing the upload process
	if chunk.GetFinal() {
		log.Printf("(%v) Final chunk received for file: '%v'", file_id, filename)
		converted_path, err := convert.ConvertFile(filename, output_path)
		if err != nil {
			log.Printf("(%v) Upload failed: %v", file_id, err)
		}
		log.Printf("(%v) Finished conversion: %v", file_id, converted_path)

		// We're done working with this file, so delete it from the map
		delete(uploads, request.GetFileName())
	}

	return &uploadPb.ChunkResponse{
		Success: true,
		Message: "Received chunk!",
	}, nil
}
