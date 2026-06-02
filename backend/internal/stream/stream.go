package stream

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"lute/internal/storage"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/minio/minio-go/v7"
)

type StreamChunk struct {
	Size     int
	Data     []byte
	Sequence int
}

type SyncSession struct {
	GuestIds  []uuid.UUID
	HostIds   uuid.UUID
	SessionId uuid.UUID
	Playing   bool
}

var SyncSessions = make([]SyncSession, 0)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

func createSyncSession(conn *websocket.Conn) {
	log.Printf("%v creating new session...", conn.RemoteAddr())

	session := SyncSession{}
	s, err := uuid.NewUUID()
	if err == nil {
		session.SessionId = s
	} else {
		log.Printf("failed to create new session for socket: %v", conn.RemoteAddr())
		return
	}
	SyncSessions = append(SyncSessions, session)
}

func checkOrigin(r *http.Request) bool {
	return true
}

func AudioStream(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %v", r.Method, r.URL.Path)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[%v] Failed to start websocket: %v", r.RemoteAddr, err)
	}
	defer conn.Close()
	log.Printf("[%v] Stream started.", conn.RemoteAddr())
	// We should check out PreparedMessage from the websocket library

	// step 1: chunk audio
	// step 2: send audio via websocket
	query := r.URL.Query()
	track := query.Get("track")
	file, err := storage.MinioClient.GetObject(context.Background(), "lute-audio", track, minio.GetObjectOptions{})
	if err != nil {
		log.Printf("[%v] Failed to open file for streaming: %v", conn.RemoteAddr(), err)
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "Failed to open file."))
		time.Sleep(10 * time.Second)
		conn.Close()
		return
	}

	for {
		messageType, r, err := conn.ReadMessage()
		if err != nil {
			log.Printf("")
			return
		}

	}
	// Create the session for synchronized play
	createSyncSession(conn)
	streamBuffer := make([]byte, 8192*5)
	sequence := 0
	for {
		chunkSize, err := file.Read(streamBuffer)
		if err == io.EOF {
			log.Printf("[%v] Stream complete!", conn.RemoteAddr())
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Stream complete."))
			time.Sleep(10 * time.Second)
			conn.Close()
			return
		}
		if err != nil {
			log.Printf("[%v] Failed to chunk stream: %v", conn.RemoteAddr(), err)
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "Failed to chunk stream."))
			time.Sleep(10 * time.Second)
			conn.Close()
			return
		}
		chunk := StreamChunk{
			Size:     chunkSize,
			Data:     streamBuffer[:chunkSize],
			Sequence: sequence,
		}
		message, err := json.Marshal(&chunk)
		if err != nil {
			log.Printf("[%v] Failed to marshal chunk: %v", conn.RemoteAddr(), err)
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "Failed to chunk stream."))
			time.Sleep(10 * time.Second)
			conn.Close()
			return
		}
		conn.WriteMessage(websocket.BinaryMessage, message)
		sequence++
	}
}
