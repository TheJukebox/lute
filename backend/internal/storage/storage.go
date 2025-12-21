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
var ctx = context.Background()

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
    log.Println("Creating buckets...")
    exists, err := MinioClient.BucketExists(ctx, "lute-audio")
    if !exists && err == nil {
        log.Println("Creating bucket 'lute-audio'...")
        err = MinioClient.MakeBucket(ctx, "lute-audio", minio.MakeBucketOptions{})
    }
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
    Artist string
    Album string
    TrackNumber int
    DiskNumber int
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
        name := body.Name
        uriName := body.UriName
        contentType := body.ContentType
        ext, _ := mimeToExtension[body.ContentType]
        path := body.UriName + ext
        artist := body.Artist
        album := body.Album
        number := body.TrackNumber
        disk := body.DiskNumber

        expiry := 10 * time.Minute

        policy := minio.NewPostPolicy()
        policy.SetBucket("lute-audio")
        policy.SetKey(path)
        policy.SetExpires(time.Now().UTC().Add(expiry))
        policy.SetContentType(contentType)
        policy.SetUserMetadata("name", name)
        policy.SetUserMetadata("artist", artist)
        policy.SetUserMetadata("album", album)

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
			Name: name,
			UriName: uriName,
			Path: path,
            Artist: artist,
            Album: album,
            TrackNumber: number,
            DiskNumber: disk,
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
            log.Printf("Failed to fetch tracks from database: %w", err)
			http.Error(w, "Failed to fetch tracks.", http.StatusInternalServerError)
			return
		}
	}
}
