package models

import (
	"github.com/google/uuid"
)

type Track struct {
	id       uuid.UUID
	name     string
	artistId uuid.UUID
	albumId  uuid.UUID
	duration int
}

type Album struct {
	id       uuid.UUID
	name     string
	artistId uuid.UUID
}

type Artist struct {
	id   uuid.UUID
	name string
}

type Playlist struct {
	id     uuid.UUID
	name   string
	public bool
}

// composite for storing tracks associated with playlists
type PlaylistTrack struct {
	id         uuid.UUID
	playlistId uuid.UUID
	trackId    uuid.UUID
	position   int
}

// accounts
type User struct {
	id       uuid.UUID
	username string
	password string
}
