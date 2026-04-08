package storage

import (
	"fmt"
	"strings"
	"log"

	"github.com/google/uuid"
)

// OnDelete types
const (
    Cascade string = "CASCADE"
    SetNull string = "SET NULL"
    Restrict string = "RESTRICT"
    NoAction string = "NO ACTION"
    SetDefault string = "SET DEFAULT"
)

// Field types
const (
	TextField string = "varchar(255)"
	IntegerField string = "int"
	IDField string = "uuid"
)

type ForeignKey struct {
	table *Table
	field *Field
}

type Field struct {
	name string
	value any
	fieldType string
	primaryKey bool
	foreignKey bool
	hasDefault bool
	foreignKeyRef ForeignKey
	nullable bool
	onDelete string
	defaultValue any
}

type Table struct {
	name string
	fields []Field
}
func (t *Table) Create() {
	queryBase := `
		CREATE TABLE IF NOT EXISTS %v
		(%v, PRIMARY KEY (%v));
	`
	var fieldStrings []string
	var primaryKeys []string
	for _, field := range t.fields {
		var fieldString string
		constraintString := ""
		if !field.nullable {
			constraintString = "NOT NULL"
		}
		defaultString := ""
		if field.hasDefault {
			switch field.fieldType {
				case TextField:
					defaultString = fmt.Sprintf("DEFAULT '%v'", field.defaultValue)
				case IntegerField:
					defaultString = fmt.Sprintf("DEFAULT %v", field.defaultValue)
				case IDField:
					defaultString = "DEFAULT gen_random_uuid()"
			}
		}
		if field.foreignKey {
			fieldString = fmt.Sprintf(
				"%v %v %v %v REFERENCES %v(%v) ON DELETE %v",
				field.name,
				field.fieldType,
				constraintString,
				defaultString,
				field.foreignKeyRef.table.name,
				field.foreignKeyRef.field.name,
				field.onDelete,
			)
		} else {
			fieldString = fmt.Sprintf(
				"%v %v %v %v",
				field.name,
				field.fieldType,
				constraintString,
				defaultString,
			)
		}
		fieldStrings = append(fieldStrings, fieldString)
		if field.primaryKey {
			primaryKeys = append(primaryKeys, field.name)
		}
	}
	query := fmt.Sprintf(
		queryBase,
		t.name,
		strings.Join(fieldStrings, ", "),
		strings.Join(primaryKeys, ", "),
	)
	_, err := pool.Exec(ctx, query)
	if err != nil {
		log.Fatalf("Failed to create table '%v': %v", t.name, err)
	} else {
		log.Printf("Created '%v'.", t.name)
	}
}

// Request and response types
type UploadRequest struct {
    Name string `json:"name"`
    UriName string `json:"uriName"`
    ContentType string `json:"contentType"`
    Artist string `json:"artist"`
    Album string `json:"album"`
    TrackNumber int `json:"trackNumber"`
    DiskNumber int `json:"diskNumber"`
}

type PresignedUploadResponse struct {
    URL string `json:"url"`
    Fields map[string]string `json:"fields"`
}

type TrackResponse struct {
    Id uuid.UUID `json:"id"`
    Title string `json:"title"`
    UriName string `json:"uriName"`
    Path string `json:"path"`
    Artist Artist `json:"artist"`
    Album Album `json:"album"`
    TrackNumber int `json:"trackNumber"`
    DiskNumber int `json:"diskNumber"`
}

type AlbumResponse struct {
    Id uuid.UUID `json:"id"`
    Title string `json:"title"`
    Artist string `json:"artist"`
}

type TracksResponse struct {
	Tracks []TrackResponse `json:"tracks"`
}

type ArtistResponse struct {
    Id uuid.UUID `json:"id"`
    Name string `json:"name"`
}

type ArtistsResponse struct {
    Artists []ArtistResponse `json:"artists"`
}

type PlaylistResponse struct {
	Id uuid.UUID `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Tracks []TrackResponse `json:"tracks"`
}

type PlaylistRequest struct {
	Type string `json:"type"`
	Id string `json:"id"` 
	Title string `json:"title"`
	Description string `json:"description"`
	Track string `json:"track_id"`
}

