package endpoints

import (
	"lute/api/handlers"
	"lute/internal/stream"
	"net/http"
)

// You can register new endpoints here and they will get
// automatically picked up.
func init() {
	// Streaming
	http.HandleFunc("/stream", stream.AudioStream)

	// Tracks
	http.HandleFunc("GET /tracks", handlers.Tracks)
	http.HandleFunc("GET /artists/{artistId}/tracks", handlers.GetArtistTracks)
	http.HandleFunc("GET /albums/{albumId}/tracks", handlers.GetAlbumTracks)
	http.HandleFunc("POST /upload", handlers.Upload)
	http.HandleFunc("OPTIONS /upload", handlers.Upload)

	// Artists
	http.HandleFunc("GET /artists", handlers.GetArtists)
	http.HandleFunc("GET /artists/{artistId}", handlers.GetArtist)

	// Albums
	http.HandleFunc("GET /albums/{albumId}", handlers.GetAlbum)

	// Playlists
	http.HandleFunc("POST /playlists", handlers.Playlists)
	http.HandleFunc("OPTIONS /playlists", handlers.Playlists)
	http.HandleFunc("GET /playlists", handlers.GetPlaylists)
	http.HandleFunc("GET /playlists/{playlistId}", handlers.GetPlaylist)
	http.HandleFunc("GET /playlists/{playlistId}/tracks", handlers.GetPlaylistTracks)
}
