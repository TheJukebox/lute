package lute_db

import (
	models "lute/internal/lute_db/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
)

func CreateTrack(conn *pgx.Conn, track models.Track) uuid.UUID {
    conn.Query(
        `
        INSERT INTO tracks(
            id,
            path,
            name,
            artistid,
            albumid,
            duration
        ) VALUES (
            $1, $2, $3, $4, $5, $6
        ) RETURNING id
        `,
        track.Id,
        track.Path,
        track.Name,
        track.ArtistId,
        track.AlbumId,
        track.Duration,
    )
    return track.Id
}

func ReadTrack(conn *pgx.Conn, uuid string) (models.Track, error) {
    var track models.Track
    err := conn.QueryRow(
        `
        SELECT 
            id,
            path,
            name,
            artistid,
            albumid,
            duration
        FROM tracks
        WHERE id = $1
        `,
        uuid,
    ).Scan(
        &track.Id,
        &track.Path,
        &track.Name,
        &track.ArtistId,
        &track.AlbumId,
        &track.Duration,
    )
    return track, err
}

func UpdateTrack(conn *pgx.Conn, track models.Track) error {
    err := conn.Exec(
        `
        UPDATE tracks SET
            path = $2,
            name = $3,
            artistid = $4
            albumid = $5,
            duration = $6
        WHERE
            id = $1
        RETURNING id
        `,
        track.Id,
        track.Path,
        track.Name,
        track.ArtistId,
        track.AlbumId,
        track.Duration,
    )
    return err
}

func DeleteTrack(conn *pgx.Conn, uuid string) error {
    err := conn.Exec(
        `DELETE FROM tracks WHERE id = $1`,
        uuid,
    )
    return err
}

