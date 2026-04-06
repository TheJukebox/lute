package storage

import (
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

var ArtistsTable = Table {
	name: "artists",
	fields: []Field {
		{ name: "id", nullable: false, primaryKey: true, fieldType: IDField, hasDefault: true },
		{ name: "name", nullable: false, fieldType: TextField, hasDefault: true, defaultValue: "Unknown" },
	},
}

var AlbumsTable = Table {
	name: "albums",
	fields: []Field {
		{ name: "id", nullable: false, primaryKey: true, fieldType: IDField, hasDefault: true },
		{ name: "title", nullable: false, fieldType: TextField, hasDefault: true, defaultValue: "Unknown" },
		{ name: "artist", nullable: false, fieldType: IDField, foreignKey: true, foreignKeyRef: ForeignKey{ table: &ArtistsTable, field: &ArtistsTable.fields[0] }, onDelete: Cascade }, 
	},
}

var TracksTable = Table {
	name: "tracks",
	fields: []Field {
		{ name: "id", nullable: false, primaryKey: true, fieldType: IDField, hasDefault: true },
		{ name: "title", nullable: false, fieldType: TextField, hasDefault: true, defaultValue: "Untitled" }, 
		{ name: "uri_name", nullable: false, fieldType: TextField }, 
		{ name: "path", nullable: false, fieldType: TextField }, 
		{ name: "artist", nullable: false, fieldType: IDField, foreignKey: true, foreignKeyRef: ForeignKey{ table: &ArtistsTable, field: &ArtistsTable.fields[0] }, onDelete: Cascade }, 
		{ name: "album", nullable: false, fieldType: IDField, foreignKey: true, foreignKeyRef: ForeignKey{ table: &AlbumsTable, field: &AlbumsTable.fields[0] }, onDelete: Cascade }, 
		{ name: "track_number", fieldType: IntegerField, nullable: false, hasDefault: true, defaultValue: 1 },
		{ name: "disk_number", fieldType: IntegerField, nullable: false, hasDefault: true, defaultValue: 1 },
	},
}

var PlaylistsTable = Table {
	name: "playlists",
	fields: []Field {
		{ name: "id", nullable: false, primaryKey: true, fieldType: IDField, hasDefault: true },
		{ name: "title", nullable: false, fieldType: TextField, hasDefault: true, defaultValue: "New Playlist" }, 
		{ name: "Description", nullable: false, fieldType: TextField, hasDefault: true, defaultValue: "" }, 
	},
}

var AlbumTracks = Table {
	name: "album_tracks",
	fields: []Field {
		{ name: "album_id", nullable: false, primaryKey: true, fieldType: IDField, foreignKey: true, foreignKeyRef: ForeignKey{ table: &AlbumsTable, field: &AlbumsTable.fields[0]}, onDelete: Cascade },
		{ name: "track_id", nullable: false, primaryKey: true, fieldType: IDField, foreignKey: true, foreignKeyRef: ForeignKey{ table: &TracksTable, field: &TracksTable.fields[0]}, onDelete: Cascade },
	},
}

var ArtistTracks = Table {
	name: "artist_tracks",
	fields: []Field {
		{ name: "artist_id", nullable: false, primaryKey: true, fieldType: IDField, foreignKey: true, foreignKeyRef: ForeignKey{ table: &ArtistsTable, field: &ArtistsTable.fields[0]}, onDelete: Cascade },
		{ name: "track_id", nullable: false, primaryKey: true, fieldType: IDField, foreignKey: true, foreignKeyRef: ForeignKey{ table: &TracksTable, field: &TracksTable.fields[0]}, onDelete: Cascade },
	},
}

var ArtistAlbums = Table {
	name: "artist_albums",
	fields: []Field {
		{ name: "artist_id", nullable: false, primaryKey: true, fieldType: IDField, foreignKey: true, foreignKeyRef: ForeignKey{ table: &ArtistsTable, field: &ArtistsTable.fields[0]}, onDelete: Cascade },
		{ name: "album_id", nullable: false, primaryKey: true, fieldType: IDField, foreignKey: true, foreignKeyRef: ForeignKey{ table: &AlbumsTable, field: &AlbumsTable.fields[0]}, onDelete: Cascade },
	},
}

var PlaylistTracks = Table {
	name: "playlist_tracks",
	fields: []Field {
		{ name: "playlist_id", nullable: false, primaryKey: true, fieldType: IDField, foreignKey: true, foreignKeyRef: ForeignKey{ table: &PlaylistsTable, field: &PlaylistsTable.fields[0]}, onDelete: Cascade },
		{ name: "track_id", nullable: false, primaryKey: true, fieldType: IDField, foreignKey: true, foreignKeyRef: ForeignKey{ table: &TracksTable, field: &TracksTable.fields[0]}, onDelete: Cascade },
		{ name: "position", nullable: false, primaryKey: true, fieldType: IntegerField, hasDefault: true, defaultValue: 1, onDelete: Cascade},
	},
}

var Tables = []Table {
    ArtistsTable,
    AlbumsTable,
    TracksTable,
    PlaylistsTable,
    ArtistAlbums,
    AlbumTracks,
    ArtistTracks,
    PlaylistTracks,
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
        table.Create()
    }
    log.Println("Created all tables!")
}

type Artist struct {
   ID uuid.UUID
   Name string
}

func (obj *Artist) Create() error {
    query := `
        INSERT INTO artists (name)
        VALUES ($1)
        RETURNING id, name;
    `
    row := pool.QueryRow(ctx, query, obj.Name)
    row.Scan(&obj.ID, &obj.Name)
    return nil
}

type Album struct {
    ID uuid.UUID
    Title string
    Artist uuid.UUID
}

func (obj *Album) Create() error {
    query := `
        INSERT INTO albums (title, artist)
        VALUES ($1, $2)
        RETURNING id, title, artist;
    `
    row := pool.QueryRow(ctx, query, obj.Title, obj.Artist)
    row.Scan(&obj.ID, &obj.Title, &obj.Artist)

    query = `
        INSERT INTO artist_albums (albums_id, artists_id)
        VALUES ($1, $2);
    `
    _ = pool.QueryRow(ctx, query, obj.ID, obj.Artist)
    return nil
}

type Track struct {
	ID uuid.UUID `json:"id"` 
	Title string `json:"title"`
	UriName string `json:"uri_name"`
	Path string `json:"path"`
	Artist uuid.UUID `json:"arist_id"` 
	Album uuid.UUID `json:"album_id"` 
	TrackNumber int `json:"track_number"`
	DiskNumber int `json:"disk_number"`
}

func (obj *Track) Create() {
	query := `
		INSERT INTO tracks (title, uri_name, path, artist, album, track_number, disk_number)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, title, uri_name, path, artist, album, track_number, disk_number;
	`
	row := pool.QueryRow(ctx, query, obj.Title, obj.UriName, obj.Path, obj.Artist, obj.Album, obj.TrackNumber, obj.DiskNumber)
	err := row.Scan(&obj.ID, &obj.Title, &obj.UriName, &obj.Path, &obj.Artist, &obj.Album, &obj.TrackNumber, &obj.DiskNumber)
	if err != nil {
		log.Printf("Failed to insert track: %v", err)
	}

    query = `
        INSERT INTO artist_tracks (track_id, artist_id)
        VALUES ($1, $2);
    `
	_, err = pool.Exec(ctx, query, obj.ID, obj.Artist)
	if err != nil {	
		log.Printf("Failed to add track '%v' for artist '%v'", obj.ID, obj.Artist)
	}
    query = `
        INSERT INTO album_tracks (track_id, album_id)
        VALUES ($1, $2);
    `
    _, err = pool.Exec(ctx, query, obj.ID, obj.Album)
	if err != nil {	
		log.Printf("Failed to add track '%v' for album '%v'", obj.ID, obj.Album)
	}
}

type Playlist struct {
	ID uuid.UUID
	Title string
	Description string
	Tracks []Track
}

func (obj *Playlist) Create() error {
	query := `
		INSERT INTO playlists (title, description)
		VALUES ($1, $2)
		RETURNING id, title, description
	`
	row := pool.QueryRow(ctx, query, obj.Title, obj.Description)
	row.Scan(&obj.ID, &obj.Title, &obj.Description)
	return nil
}

func (obj *Playlist) AddTrack(track *Track) error {
	query := `
		INSERT INTO playlist_tracks (playlist_id, track_id, track_number)
		SELECT $1, $2, COALESCE(MAX(track_number), 0) + 1
		FROM playlist_tracks
		WHERE playlists_id = $1
	`
	row := pool.QueryRow(ctx, query, obj.ID, track.ID)
	row.Scan(&obj.ID, &track.ID)
	return nil
}

// working on this concept, but i dont think its needed yet
// type Track struct {
// 	ID IDField 
// 	Name TextField 
// 	UriName TextField 
// 	Path TextField
//     Artist TextField
//     Album TextField
//     TrackNumber IntegerField
//     DiskNumber IntegerField
// }
// func (t Track) Fields() []Field {
//     return []Field {
//         t.ID,
//         t.Name,
//         t.UriName,
//         t.Path,
//         t.Artist,
//         t.Album,
//         t.TrackNumber,
//         t.DiskNumber,
//     }
// }
// func Insert(row Row, table Table) error {
//     columnString := make([]string, len(table.Fields()))
//     for i, field := range table.Fields() {
//         columnString[i] = field.Name()
//     }
// 
//     valueString := make([]string, len(row.Fields()))
//     values := make([]string, len(row.Fields()))
//     for i, field = range row.Fields() {
//         valueString[i] = fmt.Sprintf("$%v", i)
//         values[i] = field.Value()
//     }
//     valueString = strings.Join(valueString, ", ")
//     query := `
//         INSERT INTO %v (%v)
//         VALUES (%v);
//     `
//     query = fmt.Sprintf(query, table.Name(), columnString, valueString)
//     _, err := pool.QueryRow(ctx, query, values...)
//     return err 
// } 
