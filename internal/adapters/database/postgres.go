package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func ConnectPostgres() *sql.DB {
	dsn := os.Getenv("POSTGRES_DSN")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to Postgres database")
	return db
}
