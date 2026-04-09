package main

import (
	"fmt"
	"flag"
	"log"
	"net/http"

	_ "lute/api/endpoints"
	"lute/internal/storage"
	"lute/internal/config"
)

func main() {
    // Configuring the logger
    log.SetFlags(log.Ltime | log.Lshortfile | log.Lmsgprefix)
    log.SetPrefix("\033[32m[SERVER]\033[0m ")

	configPath := flag.String("config", "lute-config.yaml", "path to lute config")
	flag.Parse()

	log.Printf("Reading configuration from '%v'", *configPath)
	err := config.ReadLuteConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to read configuration from '%v': %v", *configPath, err)
	}

	err = storage.Init()
    if err != nil {
        log.Fatalf("Failed to initialise database: %v", err)
    }
	err = storage.Connect(
		fmt.Sprintf("%v:%v", config.Lute.MinioHost, config.Lute.MinioPort),
		"minioadmin",
	)
    log.Println("Initialised database.")
    if err != nil {
        log.Fatalf("Failed to connect to Minio: %v", err)
    }
    log.Println("Connected to Minio.")

    log.Printf("Starting server...")
    server := &http.Server{
        Addr: fmt.Sprintf("%v:%v", config.Lute.Host, config.Lute.Port),
    }
    go func() {
        log.Fatal(server.ListenAndServe())
    }()
    log.Printf("Server listening at %v:%v", config.Lute.Host, config.Lute.Port)
    
    // Block until we kill the process manually.
    select {}
}
