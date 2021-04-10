package app

import (
	"log"

	"github.com/imgabe/todo/pkg/store"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func OpenDatabase(path string) *sqlx.DB {
	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("error opening database: %s", err)
	}

	db.MustExec(store.CreateDatabaseStatement)

	return db
}
