package main

import (
	"fmt"
	"log"
	"lute/internal/stream"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

// Values in this struct MUST be public for us to
// unmarshal correctly!
type ServerConfig struct {
    Host string
    Port int
}

func hello(w http.ResponseWriter, req *http.Request) {
    log.Printf("Received request from: %v", req.RemoteAddr)
}

func main() {
    // Configuring the logger
    log.SetFlags(log.Ltime | log.Lshortfile | log.Lmsgprefix)
    log.SetPrefix("[SERVER] ")

    configFile, err := os.ReadFile("lute-config.yaml")
    if err != nil {
        log.Fatalf("Failed to open config file: %s", err)
    }

    log.Printf("Reading config file: lute-config.yaml")
    config := ServerConfig{}

    err = yaml.Unmarshal(configFile, &config)
    if err != nil {
        log.Fatalf("Failed to parse config file: %s", err)
    }
    
    http.HandleFunc("/hello", hello)
    http.HandleFunc("/ws", stream.WebsocketHandler)
    http.HandleFunc("/stream", stream.AudioStream)

    log.Printf("Starting server...")
    server := &http.Server{
        Addr: fmt.Sprintf("%v:%v", config.Host, config.Port),

    }
    go func() {
        log.Fatal(server.ListenAndServe())
    }()
    log.Printf("Server listening at %v:%v", config.Host, config.Port)
    
    // Block until we kill the process manually.
    select {}
}
