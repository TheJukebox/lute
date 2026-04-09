package stream

import (
	"context"
	"encoding/json"
	"bytes"
	"io"
	"log"
	"lute/internal/storage"
	"net/http"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
	"github.com/minio/minio-go/v7"
)

type StreamChunk struct {
    Size int
    Data []byte
    Sequence int
}

var upgrader = websocket.Upgrader {
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
    CheckOrigin: checkOrigin,
}

// ffprobe to ID file
func fileType(file io.Reader) {
	cmd := exec.Command(
		"ffprobe",
		"-v", "error",
		"-show_entries", "format=format_name",
		"-of", "default=noprint_wrappers=1:nokey=1",
		"-",
	)
	cmd.Stdin = file
	log.Println(cmd.Args)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Failed to determine filetype: %v", err)
	}
	log.Println(string(output))
}

// ffmpeg to grab flac header
// stream chunk type with flac header attached

func checkOrigin(r *http.Request) bool {
    return true;
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %v", r.Method, r.URL.Path)
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Failed to open websocket connection: %v", err)
        return
    }
    defer conn.Close()

    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            log.Printf("Failed to read incoming websocket message: %v", err)
            break
        }
        log.Printf("Received: %s", message)
        conn.WriteMessage(messageType, message)
        if err != nil {
            log.Printf("Failed to write a message to websocket: %v", err)
            break
        }
    }
}

func AudioStream(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %v", r.Method, r.URL.Path)
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("[%v] Failed to start stream: %v", r.RemoteAddr, err)
    }
    defer conn.Close()
    log.Printf("[%v] Stream started.", conn.RemoteAddr())
    // We should check out PreparedMessage from the websocket library

    // step 1: chunk audio
    // step 2: send audio via websocket
    query := r.URL.Query()
    track := query.Get("track")
    file, err := storage.MinioClient.GetObject(context.Background(), "lute-audio", track, minio.GetObjectOptions{});
	if err != nil {
		log.Printf("Failed to fetch file from minio: %v", err)
	}
	defer file.Close()
	var fileBuffer bytes.Buffer
	if _, err := io.Copy(&fileBuffer, file); err != nil {
		log.Printf("Failed to process buffer: %v", err)
	}
	fileTypeReader := bytes.NewReader(fileBuffer.Bytes())
	fileType(fileTypeReader)

    if err != nil {
        log.Printf("[%v] Failed to open file for streaming: %v", conn.RemoteAddr(), err) 
            conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "Failed to open file."))
            time.Sleep(10 * time.Second)
            conn.Close()
            return
    }
    streamBuffer := make([]byte, 8192*5)
    sequence := 0
	reader := bytes.NewReader(fileBuffer.Bytes())
    for {
        chunkSize, err := reader.Read(streamBuffer) 
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
            Size: chunkSize,
            Data: streamBuffer[:chunkSize],
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
