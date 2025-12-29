package storage

import (
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

var ArtistsTable = BaseTable {
    name: "artists",
    primaryKey: IDField{ name: "id", null: false, primaryKey: true },
    fields: []Field {
        IDField{ name: "id", null: false, primaryKey: true },
        TextField{ name: "name", null: false, hasDefault: true, defaultValue: "Unknown" },
    },
}

var AlbumsTable = BaseTable {
    name: "albums",
    primaryKey: IDField{ name: "id", null: false, primaryKey: true },
    fields: []Field {
        IDField{ name: "id", null: false, primaryKey: true },
        IDField{ name: "artist", null: true, foreignKey: true, foreignKeyRef: ArtistsTable }, 
        TextField{ name: "title", null: false, hasDefault: true, defaultValue: "Unknown" },
    },
}


var AlbumTracks = JunctionTable {
    name: "album_tracks",
    referenceTables: []Table {
        AlbumsTable,
        TracksTable,
    },
}

var ArtistAlbums = JunctionTable {
    name: "artist_albums",
    referenceTables: []Table {
        ArtistsTable,
        AlbumsTable,
    },
}

var ArtistTracks = JunctionTable {
    name: "artist_tracks",
    referenceTables: []Table {
        ArtistsTable,
        TracksTable,
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
    ArtistAlbums,
    AlbumTracks,
    ArtistTracks,
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
    ID uuid.UUID
    Title string
    UriName string
    Path string
    Artist uuid.UUID 
    Album uuid.UUID 
    TrackNumber int
    DiskNumber int
}

func (obj *Track) Create() error {
	query := `
		INSERT INTO tracks (title, uri_name, path, artist, album, track_number, disk_number)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, title, uri_name, path, artist, album, track_number, disk_number;
	`
	row := pool.QueryRow(ctx, query, obj.Title, obj.UriName, obj.Path, obj.Artist, obj.Album, obj.TrackNumber, obj.DiskNumber)
    row.Scan(&obj.ID, &obj.Title, &obj.UriName, &obj.Path, &obj.Artist, &obj.Album, &obj.TrackNumber, &obj.DiskNumber)

    query = `
        INSERT INTO artist_tracks (tracks_id, artists_id)
        VALUES ($1, $2);
    `
    _ = pool.QueryRow(ctx, query, obj.ID, obj.Artist)
    query = `
        INSERT INTO album_tracks (tracks_id, albums_id)
        VALUES ($1, $2);
    `
    _ = pool.QueryRow(ctx, query, obj.ID, obj.Album)
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
