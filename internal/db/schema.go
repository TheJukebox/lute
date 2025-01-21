package db

import (
	"database/sql"
	"fmt"
)

type PsqlConn struct {
	hostname string
	port     string
	database string
	user     string
	password string
}

func Connect(conn PsqlConn) (*sql.DB, error) {
	connString := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s",
		conn.hostname,
		conn.port,
		conn.database,
		conn.user,
		conn.password,
	)
	connection, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}
	return connection, nil
}
