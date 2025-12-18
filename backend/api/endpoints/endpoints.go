package endpoints

import (
	"lute/internal/storage"
	"lute/internal/stream"
	"net/http"
)

// You can register new endpoints here and they will get
// automatically picked up.
func init() {
    http.HandleFunc("/upload", storage.Upload)
    http.HandleFunc("/stream", stream.AudioStream)
    http.HandleFunc("/ws", stream.WebsocketHandler)

	http.HandleFunc("/tracks", storage.Tracks)
}
