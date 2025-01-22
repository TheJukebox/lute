package lute_db

import (
	"fmt"
	"log"
	"lute/internal/lute_db/models"
	"reflect"
	"strings"

	"github.com/jackc/pgx"
)

type PsqlConn struct {
	Hostname string
	Port     int
	Database string
	User     string
	Password string
}

func getPsqlType(t reflect.Type) string {
	if t.String() == "int" {
		return "INT"
	}
	if t.String() == "string" {
		return "VARCHAR(255)"
	}
	if t.String() == "uuid.UUID" {
		return "uuid"
	}
	if t.String() == "bool" {
		return "boolean"
	}
	return ""
}

func createTableFromModel(fields reflect.Type, name string) string {
	var columns []string
	for i := range fields.NumField() {
		field := fields.Field(i)
		t := getPsqlType(field.Type)
		columns = append(columns, fmt.Sprintf("\t%v %v", strings.ToLower(field.Name), t))
	}
	columnsString := strings.Join(columns, ",\n")
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v (\n%v,\n\tPRIMARY KEY (id)\n);", name, columnsString)
}

func Connect(config pgx.ConnConfig) (*pgx.ConnPool, error) {
	connPoolConf := pgx.ConnPoolConfig{}
	connPoolConf.ConnConfig = config
	connection, err := pgx.NewConnPool(connPoolConf)
	if err != nil {
		return nil, err
	}
	return connection, nil
}

func CreateTables(conn *pgx.ConnPool) {
	log.Print("(Sorta but not really) Migrating database...")
	users := createTableFromModel(reflect.TypeOf(models.User{}), "users")
	userConn, err := conn.Acquire()
	if err != nil {
		log.Printf("Failed db migration: %v", err)
	}
	userConn.Query(users)
	log.Print("\t Migrated 'users'...")
	userConn.Close()

	tracks := createTableFromModel(reflect.TypeOf(models.Track{}), "tracks")
	trackConn, err := conn.Acquire()
	if err != nil {
		log.Printf("Failed db migration: %v", err)
	}
	trackConn.Query(tracks)
	log.Print("\t Migrated 'tracks'...")
	trackConn.Close()

	albums := createTableFromModel(reflect.TypeOf(models.Album{}), "albums")
	albumConn, err := conn.Acquire()
	if err != nil {
		log.Printf("Failed db migration: %v", err)
	}
	albumConn.Query(albums)
	log.Print("\t Migrated 'albums'...")
	albumConn.Close()

	artists := createTableFromModel(reflect.TypeOf(models.Artist{}), "artists")
	artistConn, err := conn.Acquire()
	if err != nil {
		log.Printf("Failed db migration: %v", err)
	}
	artistConn.Query(artists)
	log.Print("\t Migrated 'artists'...")
	artistConn.Close()

	playlists := createTableFromModel(reflect.TypeOf(models.Playlist{}), "playlists")
	playlistConn, err := conn.Acquire()
	if err != nil {
		log.Printf("Failed db migration: %v", err)
	}
	playlistConn.Query(playlists)
	log.Print("\t Migrated 'playlists'...")
	defer playlistConn.Close()

	playlisttrack := createTableFromModel(reflect.TypeOf(models.PlaylistTrack{}), "playlisttrack")
	playlisttrackConn, err := conn.Acquire()
	if err != nil {
		log.Printf("Failed db migration: %v", err)
	}
	playlisttrackConn.Query(playlisttrack)
	log.Print("\t Migrated 'playlisttrack'...")
	playlisttrackConn.Close()
}
