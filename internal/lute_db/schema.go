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

func createTableFromModel(fields reflect.Type, name string, conn *pgx.Conn) {
	var columns []string
	for i := range fields.NumField() {
		field := fields.Field(i)
		t := getPsqlType(field.Type)
		columns = append(columns, fmt.Sprintf("\t%v %v", strings.ToLower(field.Name), t))
	}
	columnsString := strings.Join(columns, ",\n")
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v (\n%v,\n\tPRIMARY KEY (id)\n);", name, columnsString)
	_, err := conn.Query(query)
	if err != nil {
		log.Printf("Failed to migrate db: %v", err)
	}
	defer conn.Close()
}

func Connect(config pgx.ConnConfig) (*pgx.ConnPool, error) {
	connPoolConf := pgx.ConnPoolConfig{}
	connPoolConf.ConnConfig = config
	connPoolConf.MaxConnections = 10
	connection, err := pgx.NewConnPool(connPoolConf)
	if err != nil {
		return nil, err
	}
	return connection, nil
}

func CreateTables(conn *pgx.ConnPool) {
	log.Print("(Sorta but not really) Migrating database...")
	newConn, err := conn.Acquire()
	if err != nil {
		log.Printf("Failed db migration: %v", err)
	}
	createTableFromModel(reflect.TypeOf(models.User{}), "users", newConn)
	log.Print("\t Migrated 'users'...")

	newConn, err = conn.Acquire()
	if err != nil {
		log.Printf("Failed db migration: %v", err)
	}
	createTableFromModel(reflect.TypeOf(models.Track{}), "tracks", newConn)
	log.Print("\t Migrated 'tracks'...")

	newConn, err = conn.Acquire()
	if err != nil {
		log.Printf("Failed db migration: %v", err)
	}
	createTableFromModel(reflect.TypeOf(models.Album{}), "albums", newConn)
	log.Print("\t Migrated 'albums'...")

	newConn, err = conn.Acquire()
	if err != nil {
		log.Printf("Failed db migration: %v", err)
	}
	createTableFromModel(reflect.TypeOf(models.Artist{}), "artists", newConn)
	log.Print("\t Migrated 'artists'...")

	newConn, err = conn.Acquire()
	if err != nil {
		log.Printf("Failed db migration: %v", err)
	}
	createTableFromModel(reflect.TypeOf(models.Playlist{}), "playlists", newConn)
	log.Print("\t Migrated 'playlists'...")

	newConn, err = conn.Acquire()
	if err != nil {
		log.Printf("Failed db migration: %v", err)
	}
	createTableFromModel(reflect.TypeOf(models.PlaylistTrack{}), "playlisttrack", newConn)
	log.Print("\t Migrated 'playlisttrack'...")
}
