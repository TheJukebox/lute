package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client
var ctx = context.Background()

func Connect(
    endpoint string,
    accessKey string,
) error {
    minioClient, err := minio.New(
        endpoint,
        &minio.Options{
            Creds: credentials.NewStaticV4("minioadmin", accessKey, ""),
            Secure: false,
        },
    )
    MinioClient = minioClient
    log.Println("Creating buckets...")
    exists, err := MinioClient.BucketExists(ctx, "lute-audio")
    if !exists && err == nil {
        log.Println("Creating bucket 'lute-audio'...")
        err = MinioClient.MakeBucket(ctx, "lute-audio", minio.MakeBucketOptions{})
    }
    return err
}

func AllArtists() ([]ArtistResponse, error) {
	query := `
		SELECT id, name FROM artists;		
	`
	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Failed to gather Artists: %w", err)
	}
    var artists []Artist
    for rows.Next() {
        var artist Artist 
		err = rows.Scan(&artist.ID, &artist.Name)
		if err != nil {
			return nil, fmt.Errorf("Failed to gather Tracks: %w", err)
		}
		artists = append(artists, artist)
    }
    response := make([]ArtistResponse, len(artists))
    for i, artist := range artists {
        response[i] = ArtistResponse {
            Id: artist.ID,
            Name: artist.Name,
        }
    }
    return response, rows.Err()
}

func ArtistByName(name string) (Artist, error) {
    query := `
        SELECT id, name FROM artists WHERE name = $1;
    `
    artist := Artist{}
    err := pool.QueryRow(ctx, query, name).Scan(&artist.ID, &artist.Name)
    return artist, err
}

func ArtistByID(id string) (Artist, error) {
    query := `
        SELECT id, name FROM artists WHERE id = $1;
    `
    artist := Artist{}
    err := pool.QueryRow(ctx, query, id).Scan(&artist.ID, &artist.Name)
    return artist, err
}


func AlbumByTitle(title string, artist uuid.UUID) (Album, error) {
    query := `
        SELECT id, title, artist FROM albums WHERE title = $1 and artist = $2;
    `
    album := Album{}
    err := pool.QueryRow(ctx, query, title, artist).Scan(&album.ID, &album.Title, &album.Artist)
    return album, err
}

func AlbumByID(id string) (Album, error) {
    query := `
        SELECT id, title, artist FROM albums WHERE id = $1;
    `
    album := Album{}
    err := pool.QueryRow(ctx, query, id).Scan(&album.ID, &album.Title, &album.Artist)
    return album, err
}

func TracksByArtist(id string) ([]TrackResponse, error) {

    query := `
        SELECT id, title, uri_name, path, artist, album, track_number, disk_number
        FROM tracks WHERE artist = $1;
    `
    rows, err := pool.Query(ctx, query, id)
    if err != nil {
        return nil, fmt.Errorf("Failed to gather tracks for artist '%v'", id)
    }
    var tracks []Track
    for rows.Next() {
        var track Track
		err = rows.Scan(&track.ID, &track.Title, &track.UriName, &track.Path, &track.Artist, &track.Album, &track.TrackNumber, &track.DiskNumber)
		if err != nil {
			return nil, fmt.Errorf("Failed to gather Tracks: %w", err)
		}
		tracks = append(tracks, track)
    }
    response := make([]TrackResponse, len(tracks))
    for i, track := range tracks {
        artist, _ := ArtistByID(track.Artist.String())
        album, _ := AlbumByID(track.Album.String())
        response[i] = TrackResponse {
            Id: track.ID,
            Title: track.Title,
            UriName: track.UriName,
            Path: track.Path,
            Artist: artist,
            Album: album,
            TrackNumber: track.TrackNumber,
            DiskNumber: track.DiskNumber,
        }
    }
	return response, rows.Err()
}

func TracksByAlbum(id string) ([]TrackResponse, error) {
    query := `
        SELECT id, title, uri_name, path, artist, album, track_number, disk_number
        FROM tracks WHERE album = $1;
    `
    rows, err := pool.Query(ctx, query, id)
    if err != nil {
        return nil, fmt.Errorf("Failed to gather tracks for artist '%v'", id)
    }
    var tracks []Track
    for rows.Next() {
        var track Track
		err = rows.Scan(&track.ID, &track.Title, &track.UriName, &track.Path, &track.Artist, &track.Album, &track.TrackNumber, &track.DiskNumber)
		if err != nil {
			return nil, fmt.Errorf("Failed to gather Tracks: %w", err)
		}
		tracks = append(tracks, track)
    }
    response := make([]TrackResponse, len(tracks))
    for i, track := range tracks {
        artist, _ := ArtistByID(track.Artist.String())
        album, _ := AlbumByID(track.Album.String())
        response[i] = TrackResponse {
            Id: track.ID,
            Title: track.Title,
            UriName: track.UriName,
            Path: track.Path,
            Artist: artist,
            Album: album,
            TrackNumber: track.TrackNumber,
            DiskNumber: track.DiskNumber,
        }
    }
	return response, rows.Err()
}

func AllTracks() ([]TrackResponse, error) {
	query := `
		SELECT id, title, uri_name, path, artist, album, track_number, disk_number FROM tracks;		
	`
	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Failed to gather Tracks: %w", err)
	}
	var tracks []Track
	for rows.Next() {
		var track Track
		err = rows.Scan(&track.ID, &track.Title, &track.UriName, &track.Path, &track.Artist, &track.Album, &track.TrackNumber, &track.DiskNumber)
		if err != nil {
			return nil, fmt.Errorf("Failed to gather Tracks: %w", err)
		}
		tracks = append(tracks, track)
	}
    response := make([]TrackResponse, len(tracks))
    for i, track := range tracks {
        artist, _ := ArtistByID(track.Artist.String())
        album, _ := AlbumByID(track.Album.String())
        response[i] = TrackResponse {
            Id: track.ID,
            Title: track.Title,
            UriName: track.UriName,
            Path: track.Path,
            Artist: artist,
            Album: album,
            TrackNumber: track.TrackNumber,
            DiskNumber: track.DiskNumber,
        }
    }
	return response, rows.Err()
}

