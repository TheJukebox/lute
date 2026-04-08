package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"lute/internal/storage"
	"net/http"
	"strings"
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

func uuidFromPath(path string) (string, error) {
	id := strings.Split(path, "/")[2]
	return id, uuid.Validate(id)
}

func Tracks(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("Failed to fetch tracks from database: %v", err)
		http.Error(w, "Failed to fetch tracks.", http.StatusInternalServerError)
		return
	}
}

func GetArtists(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %v", r.Method, r.URL.Path)
	artists, err := storage.AllArtists()	
	if err != nil {
		http.Error(w, "Failed to fetch artists.", http.StatusInternalServerError)
		return
	}
	response := storage.ArtistsResponse {
		Artists: artists,
	}
	if err := writeResponse(w, response); err != nil {
		log.Printf("Failed to fetch artists from database: %v", err)
		http.Error(w, "Failed to fetch artists.", http.StatusInternalServerError)
		return
	}
}

func GetArtist(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %v", r.Method, r.URL.Path)
	artistId, err := uuidFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid UUID.", http.StatusNotFound)
		log.Printf("Invalid artist UUID provided: '%v'", artistId)
		return
	}
	artist, err := storage.ArtistByID(artistId)
	if err != nil {
		http.Error(w, "Failed to fetch artists.", http.StatusInternalServerError)
		return
	}
	response := storage.ArtistResponse {
		Id: artist.ID,
		Name: artist.Name,
	}
	if err := writeResponse(w, response); err != nil {
		log.Printf("Failed to fetch artist from database: %v", err)
		http.Error(w, "Failed to fetch artists.", http.StatusInternalServerError)
		return
	}
}

func CreatePlaylist(w http.ResponseWriter, request storage.PlaylistRequest) (storage.Playlist, error) {
	if request.Title == "" {
		http.Error(w, "title must be specified.", http.StatusBadRequest)
		return storage.Playlist{}, errors.New("Create playlist request missing title.") 
	}
	title := request.Title
	description := request.Description
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
	if err := writeResponse(w, response); err != nil {
		http.Error(w, "Failed to create playlist.", http.StatusInternalServerError)
		return storage.Playlist{}, errors.New("Failed to create playlist.")
	}
	return playlist, nil
}

func AddTrackToPlaylist(w http.ResponseWriter, request storage.PlaylistRequest) (storage.TrackResponse, error) {
	log.Print(request.Track)
	track, err := storage.TrackByID(request.Track)
	if err != nil {
		http.Error(w, "Failed to find track.", http.StatusInternalServerError)
		return track, err
	}
	storage.AddTrackToPlaylist(track.Id.String(), request.Id) 
	if err := writeResponse(w, track); err != nil {
		http.Error(w, "Failed to add track to playlist.", http.StatusInternalServerError)
		return track, err 
	}
	return track, nil
}

func GetPlaylists(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %v", r.Method, r.URL.Path)
	response, err := storage.AllPlaylists()
	if err != nil {
		log.Printf("Failed to fetch playlists: %v", err)
		http.Error(w, "Failed to  fetch playlists.", http.StatusInternalServerError)
	}
	writeResponse(w, response)
}

func GetPlaylistTracks(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %v", r.Method, r.URL.Path)
	playlistId, err := uuidFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid UUID.", http.StatusNotFound)
		log.Printf("Invalid playlist UUID provided: '%v'", playlistId)
		return
	}
	response, err := storage.TracksByPlaylistId(playlistId)
	if err != nil {
		log.Printf("Failed to fetch tracks for playlist '%v': %v", playlistId, err)
		http.Error(w, "Failed to fetch tracks for playlist", http.StatusInternalServerError)
		return
	}
	writeResponse(w, response)
}

func Playlists(w http.ResponseWriter, r *http.Request) {	
	log.Printf("%v %v", r.Method, r.URL.Path)
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
			http.Error(w, "Failed to create new playlist.", http.StatusInternalServerError)
			return
		}
		var body storage.PlaylistRequest
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			log.Printf("[%v] Failed to decode playlist request: %v", r.RemoteAddr, err)
			http.Error(w, "Failed to create new playlist.", http.StatusInternalServerError)
			return
		}

		switch body.Type {
			case "create":
				playlist, err := CreatePlaylist(w, body)
				if err != nil {
					http.Error(w, "Failed to create playlist.", http.StatusInternalServerError)
					log.Printf("[%v] Failed to create playlist: %v", r.RemoteAddr, err)
					return
				}
				log.Print(playlist.Title)
			case "add":
				track, err := AddTrackToPlaylist(w, body)
				if err != nil {
					http.Error(w, "Failed to add track to playlist", http.StatusInternalServerError)
					log.Printf("[%v] Failed to add track to playlist: %v", r.RemoteAddr, err)
					return
				}
				log.Print(track.Title)
			default:
				http.Error(w, "type must be specified: create, add", http.StatusBadRequest)
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
	log.Printf("%v %v", r.Method, r.URL.Path)
    if r.Method == http.MethodOptions {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        w.WriteHeader(http.StatusNoContent)
        return
    }

    log.Printf("[%v] Received upload request.", r.RemoteAddr)
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
