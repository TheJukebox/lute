package handlers

import (
	"context"
	"lute/internal/storage"
	"log"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/minio/minio-go/v7"
)

var mimeToExtension = map[string]string {
    "audio/mpeg": ".mp3",
    "audio/mp4": ".m4a",
    "audio/ogg": ".ogg",
    "audio/flac": ".flac",
    "audio/wav": ".wav",
}

func writeResponse(w http.ResponseWriter, response any) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func Tracks(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
        artistID := r.URL.Query().Get("artist")
        albumID := r.URL.Query().Get("album")
        var tracks []storage.TrackResponse
        var err error
        if albumID != "" {
            tracks, err = storage.TracksByAlbum(albumID)
        } else if artistID != "" {
            tracks, err = storage.TracksByArtist(artistID)
        } else {
            tracks, err = storage.AllTracks()	
            if err != nil {
                http.Error(w, "Failed to fetch tracks.", http.StatusInternalServerError)
                return
            }
        }
		response := storage.TracksResponse {
			Tracks: tracks,
		}
		if err = writeResponse(w, response); err != nil {
            log.Printf("Failed to fetch tracks from database: %w", err)
			http.Error(w, "Failed to fetch tracks.", http.StatusInternalServerError)
			return
		}
	}
}

func Artists(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
        artistID := r.URL.Query().Get("artist")
        var response any
        if artistID != "" {
            artist, err := storage.ArtistByID(artistID)
            if err != nil {
                http.Error(w, "Failed to fetch artists.", http.StatusInternalServerError)
                return
            }
            response = storage.ArtistResponse {
                Id: artist.ID,
                Name: artist.Name,
            }
        } else {
            artists, err := storage.AllArtists()	
            if err != nil {
                http.Error(w, "Failed to fetch artists.", http.StatusInternalServerError)
                return
            }
            response = storage.ArtistsResponse {
                Artists: artists,
            }
        }
        if err := writeResponse(w, response); err != nil {
            log.Printf("Failed to fetch artists from database: %w", err)
			http.Error(w, "Failed to fetch artists.", http.StatusInternalServerError)
			return
		}
	}
}


func Playlists(w http.ResponseWriter, r *http.Request) {	
    if r.Method == http.MethodOptions {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        w.WriteHeader(http.StatusNoContent)
        return
    }

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Printf("[%v] Failed to create new playlist: %v", r.RemoteAddr, err)
		}
		var body storage.CreatePlaylistRequest
		err = json.NewDecoder(r.Body).Decode(&body)
		log.Print(err)
		if err != nil || body.Title == "" {
			http.Error(w, "title must be specified.", http.StatusBadRequest)
			return
		}
		title := body.Title
		description := body.Description
		playlist := storage.Playlist {
			Title: title,
			Description: description,
		}
		playlist.Create()
		response := storage.PlaylistResponse {
			Id: playlist.ID,
			Title: playlist.Title,
			Description: playlist.Description,
		}
        if err = writeResponse(w, response); err != nil {
            http.Error(w, "Failed to generate a presigned URL.", http.StatusInternalServerError)
            return
        }
	}
}

func findExistingArtistByName(name string) storage.Artist {
    artist, err := storage.ArtistByName(name)
    if err != nil && err == pgx.ErrNoRows {
        log.Printf("No artist '%v' exists currently. Creating...", name)
        artist.Name = name
        artist.Create()
    }
    return artist
}

func findExistingAlbum(title string, artistID uuid.UUID) (storage.Album) {
    album, err := storage.AlbumByTitle(title, artistID)
    if err != nil && err == pgx.ErrNoRows {
        log.Printf("No album '%v' by '%v' exists currently. Creating...", title, artistID)
        album.Title = title
        album.Artist = artistID
        album.Create()
    }
    return album
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
        var body storage.UploadRequest
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
        artistName := body.Artist
        albumName := body.Album
        number := body.TrackNumber
        disk := body.DiskNumber

        expiry := 10 * time.Minute

        policy := minio.NewPostPolicy()
        policy.SetBucket("lute-audio")
        policy.SetKey(path)
        policy.SetExpires(time.Now().UTC().Add(expiry))
        policy.SetContentType(contentType)
        policy.SetUserMetadata("name", name)
        policy.SetUserMetadata("artist", artistName)
        policy.SetUserMetadata("album", albumName)

        presignedURL, formData, err := storage.MinioClient.PresignedPostPolicy(context.Background(), policy)
        if err != nil {
            log.Printf("[%v] Failed to generate a presigned URL: %v", r.RemoteAddr, err)
            http.Error(w, "Failed to generate a presigned URL.", http.StatusInternalServerError)
            return
        }
        response := storage.PresignedUploadResponse {
            URL: presignedURL.String(),
            Fields: formData,
        }
        if err = writeResponse(w, response); err != nil {
            http.Error(w, "Failed to generate a presigned URL.", http.StatusInternalServerError)
            return
        }

        artist := findExistingArtistByName(artistName)
        album := findExistingAlbum(albumName, artist.ID)

		track := storage.Track {
			Title: name,
			UriName: uriName,
			Path: path,
            Artist: artist.ID,
            Album: album.ID,
            TrackNumber: number,
            DiskNumber: disk,
		}
		track.Create()
    }
}
