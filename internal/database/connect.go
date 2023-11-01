package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

func (database *Database) Connect() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "admin", "admin", "admin")

	db, err := sqlx.Connect("pgx", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	database.db = db
}
