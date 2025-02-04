package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func ConnectPostgres(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Fatal("cannot ping database:", err)
		return nil, err
	}

	return db, nil
}
