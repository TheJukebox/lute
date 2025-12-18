package storage

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func Connect(
    endpoint string,
    accessKey string,
) error {
    minioClient, err := minio.New(
        endpoint,
        &minio.Options{
            Creds: credentials.NewStaticV4("minioadmin", accessKey, ""),
            Secure: false,
        },
    )
    MinioClient = minioClient
    return err
}

type UploadRequest struct {
    Name string `json:"name"`
    UriName string `json:"uriName"`
    ContentType string `json:"contentType"`
}

func Upload(w http.ResponseWriter, r *http.Request) {
    log.Printf("[%v] Received upload request.", r.RemoteAddr)
    if r.Method == http.MethodPost {
        err := r.ParseForm()
        if err != nil {
            log.Printf("[%v] Failed to parse upload request: %v", r.RemoteAddr, err)
        }
        var body UploadRequest
        err = json.NewDecoder(r.Body).Decode(&body)
        if err != nil || body.Name == "" || body.UriName == "" || body.ContentType == "" {
            w.Header().Set("Content-Type", "text/plain")
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte("name, uriName, and contentType must be specified."))
            return
        }

        expiry := 10 * time.Minute
        presignedURL, err := MinioClient.PresignedPutObject(context.Background(), "lute-audio", body.UriName, expiry)
        if err != nil {
            log.Printf("[%v] Failed to generate a presigned URL: %v", r.RemoteAddr, err)
            w.Header().Set("Content-Type", "text/plain")
            w.WriteHeader(http.StatusInternalServerError)
            w.Write([]byte("The server failed to generate a presigned URL."))
            return
        }
        w.Header().Set("Content-Type", "text/plain")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(presignedURL.String()))
    }
}
