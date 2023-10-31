package database

import (
	"encoding/json"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"wb/internal/model"
)

type Database struct {
	db *sqlx.DB
}

func (database *Database) Consume(data []byte) error {
	var order model.Order
	err := json.Unmarshal(data, &order)
	if err != nil {
		log.Println("database consume unmarshal: ", err)
		return err
	}

	database.Connect()

	err = database.AddOrder(order)
	if err != nil {
		log.Println("database consume addOrder: ", err)
		return err
	}
	return nil
}
