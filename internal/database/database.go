package database

import (
	"context"
	"encoding/json"
	"fmt"
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

func (database *Database) Connect() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "admin", "admin", "admin")

	db, err := sqlx.Connect("pgx", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected!")

	database.db = db
}

func (database *Database) AddOrder(order model.Order) error {
	tx, err := database.db.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := `
INSERT INTO orders (
    order_uid, 
    track_number,
    entry,
    locale,
    internal_signature,
    customer_id,
    delivery_service,
    shardkey,
    sm_id,
    date_created,
    oof_shard
)
VALUES (:order_uid, :track_number, :entry, :locale, :internal_signature, :customer_id, :delivery_service, :shardkey, :sm_id, :date_created, :oof_shard)
`
	_, err = tx.NamedExecContext(context.Background(), query, order)
	if err != nil {
		return err
	}

	query = `
INSERT INTO delivery (
                      order_uid, 
                      name,
                      phone,
                      zip,
                      city,
                      address,
                      region,
                      email
)

VALUES (:order_uid,:name,:phone,:zip,:city,:address,:region,:email)
`
	delivery := struct {
		OrderUID string `db:"order_uid"`
		model.Delivery
	}{
		order.OrderUID,
		order.Delivery,
	}

	_, err = tx.NamedExecContext(context.Background(), query, delivery)
	if err != nil {
		return err
	}

	query = `
INSERT INTO payment (
                     order_uid,
                     transaction,
                     request_id,
                     currency,
                     provider,
                     amount,
                     payment_dt,
                     bank,
                     delivery_cost,
                     goods_total,
                     custom_fee
)
VALUES (:order_uid, :transaction,:request_id,:currency,:provider,:amount,:payment_dt,:bank,:delivery_cost,:goods_total,:custom_fee)
`
	payment := struct {
		OrderUID string `db:"order_uid"`
		model.Payment
	}{
		order.OrderUID,
		order.Payment,
	}

	_, err = tx.NamedExecContext(context.Background(), query, payment)
	if err != nil {
		return err
	}

	query = `
INSERT INTO items (
                   chrt_id, 
                   order_uid, 
                   track_number,
				   price, 
                   rid, 
                   name, 
                   sale,
				   size, 
                   total_price, 
                   nm_id,
				   brand, 
                   status
)
VALUES (:chrt_id,:order_uid,:track_number,:price,:rid,:name,:sale,:size,:total_price,:nm_id,:brand,:status)
ON CONFLICT (chrt_id) DO NOTHING
`
	for _, item := range order.Items {

		itemWithOrderUID := struct {
			OrderUID string `db:"order_uid"`
			model.Item
		}{
			order.OrderUID,
			item,
		}

		_, err := tx.NamedExecContext(context.Background(), query, itemWithOrderUID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (database *Database) GetOrder(orderUID string) (*model.Order, error) {
	var order model.Order
	err := database.db.Get(&order, "SELECT * FROM orders WHERE order_uid=$1", orderUID)
	return &order, err
}

func (database *Database) GetAllOrders() (map[string]model.Order, error) {
	fmt.Println("GetAllOrders")
	var orders []model.Order
	err := database.db.Select(&orders, "SELECT * FROM orders")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	orderMap := make(map[string]model.Order)
	for _, order := range orders {
		orderMap[order.OrderUID] = order
	}

	fmt.Println("[init OrderUID]: ", orderMap["testid"].OrderUID)
	return orderMap, nil
}
