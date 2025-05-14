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
