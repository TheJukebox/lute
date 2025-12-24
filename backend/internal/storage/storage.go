package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func AllArtists() ([]Artist, error) {
	query := `
		SELECT id, name FROM artists;		
	`
	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Failed to gather Tracks: %w", err)
	}
	var artists []Artist
	for rows.Next() {
		var artist Artist
		err = rows.Scan(&artist.ID, &artist.Name)
		if err != nil {
			return nil, fmt.Errorf("Failed to gather Tracks: %w", err)
		}
		artists = append(artists, artist)
	}
	return artists, rows.Err()
}

func ArtistByName(name string) (Artist, error) {
    query := `
        SELECT id, name FROM artists WHERE name = $1;
    `
    artist := Artist{}
    err := pool.QueryRow(ctx, query, name).Scan(&artist.ID, &artist.Name)
    return artist, err
}

func ArtistByID(id string) (Artist, error) {
    query := `
        SELECT id, name FROM artists WHERE id = $1;
    `
    artist := Artist{}
    err := pool.QueryRow(ctx, query, id).Scan(&artist.ID, &artist.Name)
    return artist, err
}

func AlbumByTitle(title string, artist uuid.UUID) (Album, error) {
    query := `
        SELECT id, title, artist FROM albums WHERE title = $1 and artist = $2;
    `
    album := Album{}
    err := pool.QueryRow(ctx, query, title, artist).Scan(&album.ID, &album.Title, &album.Artist)
    return album, err
}

func AlbumByID(id string) (Album, error) {
    query := `
        SELECT id, title, artist FROM albums WHERE id = $1;
    `
    album := Album{}
    err := pool.QueryRow(ctx, query, id).Scan(&album.ID, &album.Title, &album.Artist)
    return album, err
}

func TracksByArtist(id string) ([]TrackResponse, error) {

    query := `
        SELECT id, title, uri_name, path, artist, album, track_number, disk_number
        FROM tracks WHERE artist = $1;
    `
    rows, err := pool.Query(ctx, query, id)
    if err != nil {
        return nil, fmt.Errorf("Failed to gather tracks for artist '%v'", id)
    }
    var tracks []Track
    for rows.Next() {
        var track Track
		err = rows.Scan(&track.ID, &track.Title, &track.UriName, &track.Path, &track.Artist, &track.Album, &track.TrackNumber, &track.DiskNumber)
		if err != nil {
			return nil, fmt.Errorf("Failed to gather Tracks: %w", err)
		}
		tracks = append(tracks, track)
    }
    response := make([]TrackResponse, len(tracks))
    for i, track := range tracks {
        artist, _ := ArtistByID(track.Artist.String())
        album, _ := AlbumByID(track.Album.String())
        response[i] = TrackResponse {
            Id: track.ID,
            Title: track.Title,
            UriName: track.UriName,
            Path: track.Path,
            Artist: artist,
            Album: album,
            TrackNumber: track.TrackNumber,
            DiskNumber: track.DiskNumber,
        }
    }
	return response, rows.Err()
}

func TracksByAlbum(id string) ([]TrackResponse, error) {

    query := `
        SELECT id, title, uri_name, path, artist, album, track_number, disk_number
        FROM tracks WHERE album = $1;
    `
    rows, err := pool.Query(ctx, query, id)
    if err != nil {
        return nil, fmt.Errorf("Failed to gather tracks for artist '%v'", id)
    }
    var tracks []Track
    for rows.Next() {
        var track Track
		err = rows.Scan(&track.ID, &track.Title, &track.UriName, &track.Path, &track.Artist, &track.Album, &track.TrackNumber, &track.DiskNumber)
		if err != nil {
			return nil, fmt.Errorf("Failed to gather Tracks: %w", err)
		}
		tracks = append(tracks, track)
    }
    response := make([]TrackResponse, len(tracks))
    for i, track := range tracks {
        artist, _ := ArtistByID(track.Artist.String())
        album, _ := AlbumByID(track.Album.String())
        response[i] = TrackResponse {
            Id: track.ID,
            Title: track.Title,
            UriName: track.UriName,
            Path: track.Path,
            Artist: artist,
            Album: album,
            TrackNumber: track.TrackNumber,
            DiskNumber: track.DiskNumber,
        }
    }
	return response, rows.Err()
}

func AllTracks() ([]TrackResponse, error) {
	query := `
		SELECT id, title, uri_name, path, artist, album, track_number, disk_number FROM tracks;		
	`
	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Failed to gather Tracks: %w", err)
	}
	var tracks []Track
	for rows.Next() {
		var track Track
		err = rows.Scan(&track.ID, &track.Title, &track.UriName, &track.Path, &track.Artist, &track.Album, &track.TrackNumber, &track.DiskNumber)
		if err != nil {
			return nil, fmt.Errorf("Failed to gather Tracks: %w", err)
		}
		tracks = append(tracks, track)
	}
    response := make([]TrackResponse, len(tracks))
    for i, track := range tracks {
        artist, _ := ArtistByID(track.Artist.String())
        album, _ := AlbumByID(track.Album.String())
        response[i] = TrackResponse {
            Id: track.ID,
            Title: track.Title,
            UriName: track.UriName,
            Path: track.Path,
            Artist: artist,
            Album: album,
            TrackNumber: track.TrackNumber,
            DiskNumber: track.DiskNumber,
        }
    }
	return response, rows.Err()
}

func findExistingArtistByName(name string) Artist {
    artist, err := ArtistByName(name)
    if err != nil && err == pgx.ErrNoRows {
        log.Printf("No artist '%v' exists currently. Creating...", name)
        artist.Name = name
        artist.Create()
    }
    return artist
}

func findExistingAlbum(title string, artistID uuid.UUID) (Album) {
    album, err := AlbumByTitle(title, artistID)
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

        artist := findExistingArtistByName(artistName)
        album := findExistingAlbum(albumName, artist.ID)

		track := Track {
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

func Tracks(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
        artistID := r.URL.Query().Get("artist")
        albumID := r.URL.Query().Get("album")
        var tracks []TrackResponse
        var err error
        if albumID != "" {
            tracks, err = TracksByAlbum(albumID)
        } else if artistID != "" {
            tracks, err = TracksByArtist(artistID)
        } else {
            tracks, err = AllTracks()	
            if err != nil {
                http.Error(w, "Failed to fetch tracks.", http.StatusInternalServerError)
                return
            }
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
