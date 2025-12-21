package storage

import (
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

type Field interface {
    Name() string
    Type() any
    Value() any
    Default() any
    Null() string 
    PrimaryKey() string
    ForeignKey(onDelete string) string
}

type IntegerField struct {
    name string
    value int 
    hasDefault bool
    defaultValue int
    null bool 
    primaryKey bool
    foreignKey bool
    foreignKeyRef Table
}
func (f IntegerField) Name() string { return f.name }
func (f IntegerField) Type() any { return "INTEGER" }
func (f IntegerField) Value() any { return f.value }
func (f IntegerField) Default() any { 
    if f.hasDefault {
       return fmt.Sprintf("DEFAULT %d", f.defaultValue)
    }
    return ""
}
func (f IntegerField) PrimaryKey() string {
    if f.primaryKey {
        return "PRIMARY KEY"
    } else {
        return ""
    }
}
func (f IntegerField) Null() string { 
    if f.null {
        return "NULL"
    } else {
        return "NOT NULL"
    }
}
func (f IntegerField) ForeignKey(onDelete string) string {
    if f.foreignKey {
        return fmt.Sprintf(
            "REFERENCES %v(%v) ON DELETE %v",
            f.foreignKeyRef.Name(),
            f.foreignKeyRef.PrimaryKey().Name(),
            onDelete,
        )
    }
    return ""
}

type TextField struct {
    name string
    value string
    hasDefault bool
    defaultValue string
    null bool 
    primaryKey bool
    foreignKey bool
    foreignKeyRef Table
}

func (f TextField) Name() string { return f.name }
func (f TextField) Type() any { return "VARCHAR" }
func (f TextField) Value() any { return f.value }
func (f TextField) Default() any { 
    if f.hasDefault {
       return "DEFAULT '" + f.defaultValue + "'" 
    }
    return ""
}
func (f TextField) PrimaryKey() string {
    if f.primaryKey {
        return "PRIMARY KEY"
    } else {
        return ""
    }
}
func (f TextField) Null() string { 
    if f.null {
        return "NULL"
    } else {
        return "NOT NULL"
    }
}
func (f TextField) ForeignKey(onDelete string) string {
    if f.foreignKey {
        return fmt.Sprintf(
            "REFERENCES %v(%v) ON DELETE %v",
            f.foreignKeyRef.Name(),
            f.foreignKeyRef.PrimaryKey().Name(),
            onDelete,
        )
    }
    return ""
}

type IDField struct {
    name string
    value uuid.UUID 
    primaryKey bool
    foreignKey bool
    foreignKeyRef Table
    null bool
}

func (f IDField) Name() string { return f.name }
func (f IDField) Type() any { return "UUID" }
func (f IDField) Value() any { return f.value }
func (f IDField) Default() any { 
    if !f.foreignKey {
        return "DEFAULT gen_random_uuid()"
    }
    return ""
}
func (f IDField) PrimaryKey() string {
    if f.primaryKey {
        return "PRIMARY KEY"
    }
    return ""
}
func (f IDField) Null() string { return "NOT NULL" }
func (f IDField) ForeignKey(onDelete string) string {
    if f.foreignKey {
        return fmt.Sprintf(
            "REFERENCES %v(%v) ON DELETE %v",
            f.foreignKeyRef.Name(),
            f.foreignKeyRef.PrimaryKey().Name(),
            onDelete, 
        )
    }
    return ""
}


type Table interface {
    Name() string 
    PrimaryKey() Field
    Fields() []Field
    Create() error
}

type JunctionTable struct {
    name string
    primaryKey Field
    referenceTables []Table
    fields []Field
}

func (t JunctionTable) Name() string { return t.name }
func (t JunctionTable) Fields() []Field { return t.fields }
func (t JunctionTable) PrimaryKey() Field { return t.primaryKey }
func (t JunctionTable) Create() error {
    queryBase := `
        CREATE TABLE IF NOT EXISTS %v
        (%v, PRIMARY KEY (%v));
    `
    referenceTables := t.referenceTables
    // we need to set up the referenceTables here
    // as sql entries
    fieldStrings := make([]string, len(referenceTables))
    pkeys := make([]string, len(referenceTables))
    for i, table := range referenceTables {
        fieldString := fmt.Sprintf(
            "%v_id %v REFERENCES %v(%v) ON DELETE CASCADE",
            table.Name(),
            table.PrimaryKey().Type(),
            table.Name(),
            table.PrimaryKey().Name(),
        )
        fieldStrings[i] = fieldString
        pkeys[i] = fmt.Sprintf("%v_id", table.Name()) 
        t.fields = append(
            t.fields, 
            IDField { name: fmt.Sprintf("%v_id", table.Name()), foreignKey: true, foreignKeyRef: table },
        )
    }
    t.primaryKey = IDField {
        name: fmt.Sprintf("pkey_%v", t.name),
    }
    query := fmt.Sprintf(
        queryBase,
        t.name,
        strings.Join(fieldStrings, ", "),
        strings.Join(pkeys, ", "),
    )
    _, err := pool.Exec(ctx, query)
    return err
}

type BaseTable struct {
    name string
    primaryKey Field
    fields []Field
}

func (t BaseTable) Name() string { return t.name }
func (t BaseTable) Fields() []Field { return t.fields }
func (t BaseTable) PrimaryKey() Field { return t.primaryKey }
func (t BaseTable) Create() error {
    queryBase := `
        CREATE TABLE IF NOT EXISTS %v
        (%v);
    `
    fields := t.Fields()
    if len(fields) == 0 {
        return fmt.Errorf("Table has no configured fields (%d)", len(fields))
    }
    fieldStrings := make([]string, len(fields))
    for i, field := range t.Fields() {
        fieldString := fmt.Sprintf(
            "%v %v %v %v %v %v",
            field.Name(),
            field.Type(),
            field.ForeignKey("CASCADE"),
            field.PrimaryKey(),
            field.Default(),
            field.Null(),
        )
        fieldStrings[i] = strings.TrimSpace(fieldString)
    }
    query := fmt.Sprintf(queryBase, t.Name(), strings.Join(fieldStrings, ", ")) 
    _, err := pool.Exec(ctx, query)
    return err
}

const (
    Cascade string = "CASCADE"
    SetNull string = "SET NULL"
    Restrict string = "RESTRICT"
    NoAction string = "NO ACTION"
    SetDefault string = "SET DEFAULT"
)

var ArtistsTable = BaseTable {
    name: "artists",
    primaryKey: IDField{ name: "id", null: false, primaryKey: true },
    fields: []Field {
        IDField{ name: "id", null: false, primaryKey: true },
    },
}

var AlbumsTable = BaseTable {
    name: "album",
    primaryKey: IDField{ name: "id", null: false, primaryKey: true },
    fields: []Field {
        IDField{ name: "id", null: false, primaryKey: true },
        IDField{ name: "artist", null: true, foreignKey: true, foreignKeyRef: ArtistsTable }, 
    },
}

var TracksTable = BaseTable {
    name: "tracks",
    primaryKey: IDField { name: "id", null: false, primaryKey: true },
    fields: []Field {
        IDField{ name: "id", null: false, primaryKey: true },
        TextField{ name: "title", null: false },
        TextField{ name: "uri_name", null: false },
        TextField{ name: "path", null: false },
        IDField{ name: "artist", null: false, foreignKey: true, foreignKeyRef: ArtistsTable },
        IDField{ name: "album", null: false, foreignKey: true, foreignKeyRef: AlbumsTable },
        IntegerField{ name: "track_number", null: false, hasDefault: true, defaultValue: 1 },
        IntegerField{ name: "disk_number", null: false, hasDefault: true, defaultValue: 1 },
    },
}

var PlaylistsTable = BaseTable {
    name: "playlists",
    primaryKey: IDField { name: "id", null: false, primaryKey: true },
    fields: []Field {
        IDField { name: "id", null: false, primaryKey: true },
        TextField { name: "title", null: false },
        TextField { name: "description", null: false, hasDefault: true, defaultValue: "A new playlist." },
    },
}

var PlaylistTracksTable = JunctionTable {
    name: "playlist_tracks",
    referenceTables: []Table {
        PlaylistsTable,
        TracksTable,
    },
}

var Tables = []Table {
    ArtistsTable,
    AlbumsTable,
    TracksTable,
    PlaylistsTable,
    PlaylistTracksTable,
}


func init() {
	log.Printf("Connecting to Postgres...")
	var err error
	pool, err = pgxpool.New(ctx, "postgres://postgres:postgres@localhost:5432/postgres")
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
    for _, table := range Tables {
        log.Printf("Creating table '%v'...", table.Name())
        err = table.Create()
        if err != nil {
            log.Fatalf("Failed to spin up tables: %w", err)
        }
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
    TrackNumber int 
    DiskNumber int 
}

func (obj Track) Create() error {
	query := `
		insert into tracks (name, uri_name, path, artist, album, track_number, disk_number)
		values ($1, $2, $3, $4, $5, $6, $7)
	`
	err := pool.QueryRow(ctx, query, obj.Name, obj.UriName, obj.Path, obj.Artist, obj.Album, obj.TrackNumber, obj.DiskNumber)
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
		err = rows.Scan(&track.ID, &track.Name, &track.UriName, &track.Path, &track.Artist, &track.Album, &track.TrackNumber, &track.DiskNumber)
		if err != nil {
			return nil, fmt.Errorf("Failed to gather Tracks: %w", err)
		}
		tracks = append(tracks, track)
	}
	return tracks, rows.Err()
}
