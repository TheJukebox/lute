package handlers

import (
	"lute/internal/storage"
	"log"
	"encoding/json"
	"net/http"
)

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
