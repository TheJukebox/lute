package endpoints

import (
	"lute/internal/stream"
	"lute/api/handlers"
	"net/http"
)

// You can register new endpoints here and they will get
// automatically picked up.
func init() {
	// Streaming
    http.HandleFunc("/stream", stream.AudioStream)
    http.HandleFunc("/ws", stream.WebsocketHandler)

	// Tracks
	http.HandleFunc("GET /tracks", handlers.Tracks)
    http.HandleFunc("POST /upload", handlers.Upload)
    http.HandleFunc("OPTIONS /upload", handlers.Upload)

	// Artists
	http.HandleFunc("GET /artists", handlers.Artists)

	// Playlists
	http.HandleFunc("POST /playlists", handlers.Playlists)
	http.HandleFunc("OPTIONS /playlists", handlers.Playlists)
	http.HandleFunc("GET /playlists", handlers.GetPlaylists)
	http.HandleFunc("GET /playlists/{playlistId}", handlers.GetPlaylistTracks)
}
