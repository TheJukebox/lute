package models

import (
	"github.com/google/uuid"
)

type Track struct {
	Id       uuid.UUID
	Path     string
	Name     string
	ArtistId uuid.UUID
	AlbumId  uuid.UUID
	Duration int
}

type Album struct {
	Id       uuid.UUID
	Name     string
	ArtistId uuid.UUID
}

type Artist struct {
	Id   uuid.UUID
	Name string
}

type Playlist struct {
	Id     uuid.UUID
	UserId uuid.UUID
	Name   string
	Public bool
}

// composite for storing tracks associated with playlists
type PlaylistTrack struct {
	Id         uuid.UUID
	PlaylistId uuid.UUID
	TrackId    uuid.UUID
	Position   int
}

// accounts
type User struct {
	Id       uuid.UUID
	Username string
	Password string
}
