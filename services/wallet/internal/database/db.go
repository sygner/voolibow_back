package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func ConnectToPostgres(dbHost string, dbPort int, dbName string, dbUsername string, dbPassword string) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		dbHost, dbPort, dbName, dbUsername, dbPassword)

	for {
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			return nil, err
		}
		err = db.Ping()
		if err == nil {
			return db, nil
		} else {
			return nil, err
		}
	}
}
