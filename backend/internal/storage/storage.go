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

var mimeToExtension = map[string]string {
    "audio/mpeg": ".mp3",
    "audio/mp4": ".m4a",
    "audio/ogg": ".ogg",
    "audio/flac": ".flac",
    "audio/wav": ".wav",
}

type UploadRequest struct {
    Name string `json:"name"`
    UriName string `json:"uriName"`
    ContentType string `json:"contentType"`
}

type PresignedUploadResponse struct {
    URL string `json:"url"`
    Fields map[string]string `json:"fields"`
}

type TracksResponse struct {
	Tracks []Track `json:"tracks"`
}

func Upload(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodOptions {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        w.WriteHeader(http.StatusNoContent)
        return
    }

    log.Printf("[%v] Received upload request.", r.RemoteAddr)
    if r.Method == http.MethodPost {
        err := r.ParseForm()
        if err != nil {
            log.Printf("[%v] Failed to parse upload request: %v", r.RemoteAddr, err)
        }
        var body UploadRequest
        err = json.NewDecoder(r.Body).Decode(&body)
        if err != nil || body.Name == "" || body.UriName == "" || body.ContentType == "" {
            http.Error(w, "name, uriName, and contentType must be specified.", http.StatusBadRequest)
            return
        }


        ext, _ := mimeToExtension[body.ContentType]
        filename := body.Name + ext 
        expiry := 10 * time.Minute

        policy := minio.NewPostPolicy()
        policy.SetBucket("lute-audio")
        policy.SetKey(filename)
        policy.SetExpires(time.Now().UTC().Add(expiry))
        policy.SetContentType(body.ContentType)
        policy.SetUserMetadata("name", body.Name)

        presignedURL, formData, err := MinioClient.PresignedPostPolicy(context.Background(), policy)
        if err != nil {
            log.Printf("[%v] Failed to generate a presigned URL: %v", r.RemoteAddr, err)
            http.Error(w, "Failed to generate a presigned URL.", http.StatusInternalServerError)
            return
        }
        response := PresignedUploadResponse {
            URL: presignedURL.String(),
            Fields: formData,
        }
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        if err = json.NewEncoder(w).Encode(response); err != nil {
            http.Error(w, "Failed to generate a presigned URL.", http.StatusInternalServerError)
            return
        }

		track := Track {
			Name: body.Name,
			UriName: body.UriName,
			Path: body.UriName + ext,
		}
		track.Create()
    }
}

func Tracks(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tracks, err := AllTracks()	
		if err != nil {
			http.Error(w, "Failed to fetch tracks.", http.StatusInternalServerError)
			return
		}
		response := TracksResponse {
			Tracks: tracks,
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to fetch tracks.", http.StatusInternalServerError)
			return
		}
	}
}
