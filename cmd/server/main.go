package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"flag"
	// how to handle cmdline flags
	// https://gobyexample.com/command-line-flags

	streamPb "lute/gen/stream"
	uploadPb "lute/gen/upload"
	api "lute/internal/lute_api"
	db "lute/internal/lute_db"
	mw "lute/internal/middleware"

	"github.com/jackc/pgx"
	"google.golang.org/grpc"
)

func debugSetup() {
	log.Printf("Creating debug folders...")
	err := os.Mkdir("uploads/raw", 0700)
	if err != nil {
		log.Printf("Failed to create uploads/raw: %v", err)
	}
	err = os.Mkdir("uploads/converted", 0700)
	if err != nil {
		log.Printf("Failed to create uploads/converted: %v", err)
	}
	log.Printf("Done!")
}

func main() {
	// config file
	type LuteConfig struct {
		Lute struct {
			Host     string `json:"host"`
			GrpcPort int    `json:"grpc"`
			HttpPort int    `json:"http"`
		} `json:"lute"`
		Postgres struct {
			PgHost string `json:"host"`
			PgPort int    `json:"port"`
		} `json:"postgres"`
		Uploads string `json:"uploads"`
		Debug   bool   `json:"debug"`
	}

	configPath := flag.String("config", "lute.config.json", "A path to a Lute configuration in JSON format.")
	flag.Parse()
	config, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}
	var confData LuteConfig
	json.Unmarshal(config, &confData)

	// args
	debug := flag.Bool("debug", confData.Debug, "Run Lute in debug mode.")
	host := flag.String("host", confData.Lute.Host, "The hostname or address that the Lute backend should listen on.")
	grpcPort := flag.Int("grpc", confData.Lute.GrpcPort, "The port that Lute should use for gRPC requests.")
	grpcAddr := fmt.Sprintf("%s:%d", *host, *grpcPort)
	httpPort := flag.Int("http", confData.Lute.HttpPort, "The port that Lute should use for HTTP requests.")
	httpAddr := fmt.Sprintf("%s:%d", *host, *httpPort)

	pgHost := flag.String("pg", confData.Postgres.PgHost, "The hostname or address of the PostgreSQL database.")
	pgPort := flag.Int("pg-port", confData.Postgres.PgPort, "The port of the PostgreSQL database.")

	flag.Parse()

	// setup folders
	if *debug {
		log.Printf("Starting Lute in debug mode...")
		debugSetup()
	}

	// db
	dbConfig := &pgx.ConnConfig{}
	dbConfig.Host = *pgHost
	dbConfig.Port = uint16(*pgPort)
	dbConfig.Database = "lute"
	dbConfig.User = "postgres"
	dbConfig.Password = "postgres"
	log.Printf("Connecting to PostgreSQL at %v:%v...", dbConfig.Host, dbConfig.Port)
	dbConnection, err := db.Connect(*dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Print("Success!")
	db.CreateTables(dbConnection)

	// start gRPC server
	log.Printf("Starting gRPC server: %s...", grpcAddr)
	listener, _ := net.Listen("tcp", grpcAddr)
	grpcNative := grpc.NewServer()
	uploadPb.RegisterUploadServer(grpcNative, &api.UploadService{})
	streamPb.RegisterAudioStreamServer(grpcNative, &api.StreamService{})
	go grpcNative.Serve(listener)
	log.Printf("Success! gRPC server listening at %v", listener.Addr())

	grpcWeb := grpc.NewServer()
	uploadPb.RegisterUploadServer(grpcWeb, &api.UploadService{})
	streamPb.RegisterAudioStreamServer(grpcWeb, &api.StreamService{})

	// standup client for HTTP/1.1 to HTTP/2
	client, err := mw.CreateGrpcClient()
	if err != nil {
		log.Fatalf("Failed to standup internal gRPC client...")
	}
	defer client.Close()
	audioClient := streamPb.NewAudioStreamClient(client)
	log.Println("Started audio streaming client for HTTP/1.1 to HTTP/2")

	// Register middleware
	mux := http.NewServeMux()
	mux.Handle("/upload.Upload/UploadFile", grpcWeb)
	mux.Handle("/stream.AudioStream/StreamAudio", grpcWeb)
	middleware := mw.GrpcWebParseMiddleware(grpcWeb, mux, audioClient)
	middleware = mw.CorsMiddleware(middleware)

	log.Printf("Starting HTTP server: %s...", httpAddr)
	server := &http.Server{
		Addr:    httpAddr,
		Handler: middleware,
	}

	log.Printf("Success! Server listening at %v", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Failed to serve: %v", err)
	}
}
