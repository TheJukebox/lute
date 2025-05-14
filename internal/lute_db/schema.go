package lute_db

import (
	"fmt"
	"log"
	"lute/internal/lute_db/models"
	"reflect"
	"strings"

	"github.com/jackc/pgx"
)

var DBConnPool *pgx.ConnPool


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

func createTableFromModel(fields reflect.Type, name string, pool *pgx.ConnPool) {
	conn, err := pool.Acquire()
	if err != nil {
		log.Printf("Failed to acquire a database connection from pool: %v", err)
	}

	var columns []string
	for i := range fields.NumField() {
		field := fields.Field(i)
		t := getPsqlType(field.Type)
		columns = append(columns, fmt.Sprintf("\t%v %v", strings.ToLower(field.Name), t))
	}
	columnsString := strings.Join(columns, ",\n")
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v (\n%v,\n\tPRIMARY KEY (id)\n);", name, columnsString)
	_, err = conn.Query(query)
	if err != nil {
		log.Printf("Failed to migrate db: %v", err)
	}
	defer conn.Close()
	log.Printf("[DATABASE] Migrated '%s'...", name)
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

func CreateTables(pool *pgx.ConnPool) {
	log.Print("(Sorta but not really) Migrating database...")

	type Table struct {
		t    reflect.Type
		name string
	}

	models := []Table{
		{reflect.TypeOf(models.User{}), "users"},
		{reflect.TypeOf(models.Track{}), "tracks"},
		{reflect.TypeOf(models.Album{}), "albums"},
		{reflect.TypeOf(models.Artist{}), "artists"},
		{reflect.TypeOf(models.Playlist{}), "playlists"},
		{reflect.TypeOf(models.PlaylistTrack{}), "playlisttrack"},
	}

	for _, model := range models {
		createTableFromModel(model.t, model.name, pool)
	}
}
