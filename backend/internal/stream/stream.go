package stream

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"lute/internal/storage"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/minio/minio-go/v7"
)

type StreamChunk struct {
    Size int
    Data []byte
    Sequence int
}

// 272 bit stream info block
type flacStreamInfo struct {
	minBlockSize uint16 // 16 bit
	maxBlockSize uint16 // 16 bit
	minFrameSize uint32 // 24 bit 
	maxFrameSize uint32 // 24 bit
	sampleRate uint32 // 20bit
	channels uint8 // 3 bit (- 1)
	bitDepth uint8 // 5 bit (- 1)
	totalSamples uint64 // 36 bits
	hash [16]byte // 128 bits
}

func (s *flacStreamInfo) Unmarshal(b []byte) {
	s.minBlockSize = uint16(b[0]) << 8 | uint16(b[1])
	s.maxBlockSize = uint16(b[2]) << 8 | uint16(b[3])
	s.minFrameSize = uint32(b[4]) << 16 | uint32(b[5]) << 8 | uint32(b[6])
	s.maxFrameSize = uint32(b[7]) << 16 | uint32(b[8]) << 8 | uint32(b[9])

	var x uint64
	for i := 10; i < 18; i ++ {
		x = (x << 8) | uint64(b[i])
	}

	s.sampleRate = uint32(x >> 44)
	s.channels = uint8((x >> 41) & 0x07) + 1
	s.bitDepth = uint8((x >> 36) & 0x1F) + 1
	s.totalSamples = x & 0xFFFFFFFFF
	copy(s.hash[:], b[18:34])
}

var upgrader = websocket.Upgrader {
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
    CheckOrigin: checkOrigin,
}

func fileType(file io.Reader) string {
	cmd := exec.Command(
		"ffprobe",
		"-v", "error",
		"-show_entries", "format=format_name",
		"-of", "default=noprint_wrappers=1:nokey=1",
		"-",
	)
	cmd.Stdin = file
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Failed to determine filetype: %v", err)
	}
	return strings.TrimSpace(string(output))
}

func getFLACStreamInfo(buf []byte) ([]byte, error) {
	marker := string(buf[:4])
	if marker != "fLaC" {
		return nil, fmt.Errorf("wanted 'flaC' in header but got '%v'", marker)
	}
	metadata := buf[4:8]
	length := uint32(metadata[1]) << 16 | uint32(metadata[2]) << 8 | uint32(metadata[3])
	streamInfo := buf[8:8+length]
	return streamInfo, nil
}

func calcCodedNumberLength(buf byte) int {
	if buf & 0x80 == 0 {
		return 1
	} else if buf & 0xE0 == 0xC0 {
		return 2
	} else if buf & 0xF0 == 0xE0 {
		return 3
	} else if buf & 0xF8 == 0xF0 {
		return 4
	} else if buf & 0xFC == 0xF8 {
		return 5
	} else if buf & 0xFE == 0xFC {
		return 6
	} else {
		return 7
	}
}

type FLACFrame struct {
	sync uint16
	strategy uint16 
	blockSize int
	sampleRate int
	channels int
	bitDepth int
}

// we need to get frame headers so we can
// chunk up a flac stream:
// https://www.rfc-editor.org/rfc/rfc9639.html#name-frame-structure
func getFLACFrames(buf []byte) {	
	// candidates are the indices that hold potential flac frames
	candidates := make([]int, 0)
	for i := 0; i < (len(buf) - 2); i++ {
		// only 15 bits of this are technically the sync, but we can capture it
		// as 0x7FFC, per the spec
		candidate := uint16(buf[i]) << 8 | uint16(buf[i + 1])
		if (candidate & 0xFFFE) == 0xFFF8 {
			candidates = append(candidates, i)
		}
	}
	// frames are indices we've confirmed are flac frames
	frames := make([]int, 0)
	for _, candidate := range candidates {
		// blocking strategy bit
		frame := FLACFrame{}
		frame.sync = uint16(buf[candidate]) << 8 | uint16(buf[candidate + 1])
		frame.strategy = (frame.sync >> 1) & 0x01
		var candidateStart uint16 
		if frame.strategy == 0 {
			candidateStart = 0xFFF8	
		} else {
			candidateStart = 0xFFF9
		}
		// ignore false positives
		if (frame.sync & candidateStart) != candidateStart {
			continue
		}
		// filter based on the coded number
		offset := candidate + 4
		if offset + 4 >= len(buf) { continue }
		codedNumberLength := calcCodedNumberLength(buf[candidate + 4])
		// validate coded number - fifth byte
		validCodedNumber := true
		for i := 1; i < codedNumberLength; i++ {
			if buf[offset + i] & 0xC0 != 0x80 {
				validCodedNumber = false
				break
			}
		}
		if !validCodedNumber {
			continue
		}
		offset += codedNumberLength

		// block size bits - first 4 bits of 3rd byte
		blockSizeBits := uint8(buf[candidate + 2]) >> 4
		if blockSizeBits == 0b0000 {
			continue
		} else if blockSizeBits == 0b0001 {
			frame.blockSize = 192	
		} else if blockSizeBits >= 0b0010 && blockSizeBits <= 0b0101 {
			frame.blockSize = 144 * (1 << blockSizeBits)		
		} else if blockSizeBits == 0b0110 {
			frame.blockSize = int(buf[offset]) + 1
			offset++
		} else if blockSizeBits == 0b0111 {
			frame.blockSize = int(uint16(buf[offset]) << 8 | uint16(buf[offset + 1])) + 1
			offset += 2
		} else if blockSizeBits >= 0b1000 && blockSizeBits <= 0b1111 {
			frame.blockSize = 1 << blockSizeBits
		} else {
			continue
		}
		// last 4 bits of 3rd byte - we can just mask the first 4
		sampleRateBits := uint8(buf[candidate + 2]) & 0x0F
		switch(sampleRateBits) {
			case 0b0000: frame.sampleRate = -1 // sampleRate is in the stream info
			case 0b0001: frame.sampleRate = 88200
			case 0b0010: frame.sampleRate = 176400 
			case 0b0011: frame.sampleRate = 12000 
			case 0b0100: frame.sampleRate = 8000 
			case 0b0101: frame.sampleRate = 16000 
			case 0b0110: frame.sampleRate = 22050
			case 0b0111: frame.sampleRate = 24000
			case 0b1000: frame.sampleRate = 32000
			case 0b1001: frame.sampleRate = 44100 
			case 0b1010: frame.sampleRate = 48000 
			case 0b1011: frame.sampleRate = 96000 
			case 0b1100: frame.sampleRate = int(buf[offset]) * 1000; offset++;
			case 0b1101: frame.sampleRate = int(uint16(buf[offset]) << 8 | uint16(buf[offset + 1])); offset += 2;
			case 0b1110: frame.sampleRate = int(uint16(buf[offset]) << 8 | uint16(buf[offset + 1])) * 10; offset += 2;
			case 0b1111: continue; // forbidden
		}
		// channel bits - first 4 bits of 4th byte
		channelBits := uint8(buf[candidate + 3]) >> 4
		switch(channelBits) {
			case 0b0000: frame.channels = 1
			case 0b0001: frame.channels = 2
			case 0b0010: frame.channels = 3
			case 0b0011: frame.channels = 4
			case 0b0100: frame.channels = 5
			case 0b0101: frame.channels = 6
			case 0b0110: frame.channels = 7
			case 0b0111: frame.channels = 8
			case 0b1000: frame.channels = 2
			case 0b1001: frame.channels = 2
			case 0b1010: frame.channels = 2
			case 0b1011, 0b1100, 0b1101, 0b1110, 0b1111: continue; // reserved
		}
		// bit depth bits - next 3 bits
		bitDepthBits := uint8(buf[candidate + 3] >> 1) & 0x07
		switch(bitDepthBits) {
			case 0b000: frame.bitDepth = -1 // bitDepth in the stream info
			case 0b001: frame.bitDepth = 8
			case 0b010: frame.bitDepth = 12
			case 0b011: continue; // reserved
			case 0b100: frame.bitDepth = 16
			case 0b101: frame.bitDepth = 20
			case 0b110: frame.bitDepth = 24
			case 0b111: frame.bitDepth = 32
		}

		log.Println(frame)
		frames = append(frames, candidate)
	}
	log.Printf("got %v frames", len(frames))
}

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
		err = conn.WriteMessage(messageType, message)
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
	reader := bytes.NewReader(fileBuffer.Bytes())

	fileType := fileType(reader)
	reader.Seek(0, io.SeekStart)
	if fileType == "flac"{
		s, err := getFLACStreamInfo(fileBuffer.Bytes())
		if err != nil {
			log.Printf("failed to read stream info: %v", err)
		}
		streamInfo := flacStreamInfo{}		
		streamInfo.Unmarshal(s)
		getFLACFrames(fileBuffer.Bytes())
		log.Println(streamInfo)
	}

    if err != nil {
        log.Printf("[%v] Failed to open file for streaming: %v", conn.RemoteAddr(), err) 
            conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "Failed to open file."))
            time.Sleep(10 * time.Second)
            conn.Close()
            return
    }
    streamBuffer := make([]byte, 8192*5)
    sequence := 0
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
