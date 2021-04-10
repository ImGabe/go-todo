package todo

import (
	"log"

	"github.com/imgabe/todo/pkg/app"
)

func main() {
	db := app.OpenDatabase(":memory:")
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("error closing database: %s", err)
		}
	}()

	if err := db.Ping(); err != nil {
		log.Fatalf("error connecting to database: %s", err)
	}
}
