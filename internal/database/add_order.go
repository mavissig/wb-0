package database

import (
	"context"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jmoiron/sqlx"
	"log"
	"wb/internal/model"
)

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
VALUES (:chrt_id,:track_number,:price,:rid,:name,:sale,:size,:total_price,:nm_id,:brand,:status)
ON CONFLICT (chrt_id) DO NOTHING
`
	fmt.Println("[Items count]: ", len(order.Items))

	for ind, item := range order.Items {

		log.Printf("ind: %d | chrt_id: %d\n", ind, item.ChrtID)
		itemWithOrderUID := struct {
			OrderUID string `db:"order_uid"`
			model.Item
		}{
			order.OrderUID,
			item,
		}

		_, err := tx.NamedExecContext(context.Background(), query, itemWithOrderUID)
		if err != nil {
			fmt.Println("[ERROR Item]: ", err)
			return err
		}
	}

	query = `
    INSERT INTO OrderItems (
        order_uid,
        chrt_id
    )
    VALUES (:order_uid, :chrt_id)
    `
	for _, item := range order.Items {
		orderItem := struct {
			OrderUID string `db:"order_uid"`
			ChrtID   int    `db:"chrt_id"`
		}{
			order.OrderUID,
			item.ChrtID,
		}

		_, err = tx.NamedExecContext(context.Background(), query, orderItem)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
