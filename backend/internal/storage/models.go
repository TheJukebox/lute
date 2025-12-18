package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool
var ctx = context.Background()

func init() {
	log.Printf("Connecting to Postgres...")
	var err error
	pool, err = pgxpool.New(ctx, "postgres://postgres:postgres@postgres:5432/postgres")
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
		return
	}
	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
		return
	}
	log.Println("Connected to Postgres!")
}

type Track struct {
	ID uuid.UUID
	Name string
	UriName string
	Path string
}

func (obj Track) Create() error {
	query := `
		INSERT INTO tracks (name, uri_name, path)
		VALUES ($1, $2, $3)
	`
	err := pool.QueryRow(ctx, query, obj.Name, obj.UriName, obj.Path)
	if err != nil {
		return fmt.Errorf("Failed to create Track object: %w", err)
	}
	return nil
}

func AllTracks() ([]Track, error) {
	query := `
		SELECT * FROM tracks;		
	`
	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Failed to gather Tracks: %w", err)
	}
	var tracks []Track
	for rows.Next() {
		var track Track
		err = rows.Scan(&track.ID, &track.Name, &track.UriName, &track.Path)
		if err != nil {
			return nil, fmt.Errorf("Failed to gather Tracks: %w", err)
		}
		tracks = append(tracks, track)
	}
	return tracks, rows.Err()
}
