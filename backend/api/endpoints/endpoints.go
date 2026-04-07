package endpoints

import (
	"lute/internal/stream"
	"lute/api/handlers"
	"net/http"
)

// You can register new endpoints here and they will get
// automatically picked up.
func init() {
    http.HandleFunc("/upload", handlers.Upload)
    http.HandleFunc("/stream", stream.AudioStream)
    http.HandleFunc("/ws", stream.WebsocketHandler)

	http.HandleFunc("/tracks", handlers.Tracks)
	http.HandleFunc("/artists", handlers.Artists)
	http.HandleFunc("POST /playlists", handlers.Playlists)
	http.HandleFunc("GET /playlists", handlers.GetPlaylists)
	http.HandleFunc("GET /playlists/{playlistId}", handlers.GetPlaylistTracks)
}
