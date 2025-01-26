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

	config, err := os.ReadFile("lute.config.json")
	if err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}
	var confData LuteConfig
	json.Unmarshal(config, &confData)

	// args
	debug := flag.Bool("debug", confData.Debug, "Run Lute in debug mode. WARNING: This will log secrets!")
	uploadPath := flag.String("uploads", confData.Uploads, "The path that Lute should store uploaded files in.")
	host := flag.String("host", confData.Lute.Host, "The hostname or address that the Lute backend should listen on.")
	grpcPort := flag.Int("grpc", confData.Lute.GrpcPort, "The port that Lute should use for gRPC requests.")
	grpcAddr := fmt.Sprintf("%s:%d", *host, *grpcPort)
	httpPort := flag.Int("http", confData.Lute.HttpPort, "The port that Lute should use for HTTP requests.")
	httpAddr := fmt.Sprintf("%s:%d", *host, *httpPort)

	pgHost := flag.String("pg", confData.Postgres.PgHost, "The hostname or address of the PostgreSQL database.")
	pgPort := flag.Int("pg-port", confData.Postgres.PgPort, "The port of the PostgreSQL database.")

	flag.Parse()

	rawPath := fmt.Sprintf("%s/raw", *uploadPath)
	convPath := fmt.Sprintf("%s/converted", *uploadPath)
	err = os.Mkdir(*uploadPath, 0700)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("Failed to create upload path: %v", err)
	}
	err = os.Mkdir(rawPath, 0700)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("Failed to create upload path: %v", err)
	}
	err = os.Mkdir(convPath, 0700)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("Failed to create upload path: %v", err)
	}

	// db
	dbConfig := &pgx.ConnConfig{}
	dbConfig.Host = *pgHost
	dbConfig.Port = uint16(*pgPort)
	dbConfig.Database = "lute"
	dbConfig.User = "postgres"
	dbConfig.Password = "postgres"

	if *debug {
		log.Printf(
			"Postgres Config:\n\thost: %v\n\tport: %v\n\tdb: %v\n\tuser: %v\n\tpassword: %v",
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.Database,
			dbConfig.User,
			dbConfig.Password,
		)
	}

	log.Printf("Connecting to PostgreSQL at %v:%v...", dbConfig.Host, dbConfig.Port)
	dbConnection, err := db.Connect(*dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Print("Success!")
	db.CreateTables(dbConnection)

	// start gRPC server
	log.Printf("Starting gRPC server on %s...", grpcAddr)
	listener, _ := net.Listen("tcp", grpcAddr)
	grpcNative := grpc.NewServer()
	uploadPb.RegisterUploadServer(grpcNative, &api.UploadService{})
	streamPb.RegisterAudioStreamServer(grpcNative, &api.StreamService{})
	go grpcNative.Serve(listener)
	log.Printf("Success! gRPC server listening on %v", listener.Addr())

	grpcWeb := grpc.NewServer()
	uploadPb.RegisterUploadServer(grpcWeb, &api.UploadService{Path: *uploadPath})
	streamPb.RegisterAudioStreamServer(grpcWeb, &api.StreamService{Path: *uploadPath})

	// standup client for HTTP/1.1 to HTTP/2
	client, err := mw.CreateGrpcClient()
	if err != nil {
		log.Fatalf("Failed to standup HTTP/1.1 to HTTP/2 middleware...")
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

	log.Printf("Starting HTTP server on %s...", httpAddr)
	server := &http.Server{
		Addr:    httpAddr,
		Handler: middleware,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("Failed to serve on %v: %v", httpAddr, err)
		}
		defer server.Close()
	}()

	log.Printf("Success! Listening on %v", httpAddr)
	select {}
}
