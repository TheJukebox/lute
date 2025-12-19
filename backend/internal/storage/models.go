package storage

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

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
    log.Println("Creating tables...")
    tracks := `
        CREATE TABLE IF NOT EXISTS tracks
        (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            name VARCHAR NOT NULL,
            uri_name VARCHAR NOT NULL,
            path VARCHAR NOT NULL,
            artist VARCHAR NOT NULL,
            album VARCHAR NOT NULL,
            number INTEGER NOT NULL,
            disk INTEGER NOT NULL
        )
    `
    _, err = pool.Query(ctx, tracks)
    if err != nil {
        log.Fatalf("Failed to create tables: %w", err)
    }
    log.Println("Created tables!")
}

type Track struct {
	ID uuid.UUID
	Name string
	UriName string
	Path string
    Artist string
    Album string
    Number int
    Disk int
}

func (obj Track) Create() error {
	query := `
		INSERT INTO tracks (name, uri_name, path, artist, album, number, disk)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	err := pool.QueryRow(ctx, query, obj.Name, obj.UriName, obj.Path, obj.Artist, obj.Album, obj.Number, obj.Disk)
	if err != nil {
		return fmt.Errorf("Failed to create Track object: %w", err)
	}
	return nil
}

func AllTracks() ([]Track, error) {
	query := `
		SELECT id, name, uri_name, path, artist, album, number, disk FROM tracks;		
	`
	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Failed to gather Tracks: %w", err)
	}
	var tracks []Track
	for rows.Next() {
		var track Track
		err = rows.Scan(&track.ID, &track.Name, &track.UriName, &track.Path, &track.Artist, &track.Album, &track.Number, &track.Disk)
		if err != nil {
			return nil, fmt.Errorf("Failed to gather Tracks: %w", err)
		}
		tracks = append(tracks, track)
	}
	return tracks, rows.Err()
}
