package database

import (
	"context"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jmoiron/sqlx"
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
ON CONFLICT (order_uid) DO UPDATE SET 
	track_number = excluded.track_number,
    entry = excluded.entry,
    locale = excluded.locale,
    internal_signature = excluded.internal_signature,
    customer_id = excluded.customer_id,
    delivery_service = excluded.delivery_service,
    shardkey = excluded.shardkey,
    sm_id = excluded.sm_id,
    date_created = excluded.date_created,
    oof_shard = excluded.oof_shard
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
ON CONFLICT (order_uid) DO UPDATE SET
	name = excluded.name,
	phone = excluded.phone,
	zip = excluded.zip,
	city = excluded.city,
	address = excluded.address,
	region = excluded.region,
	email = excluded.email
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
ON CONFLICT (transaction) DO UPDATE SET 
	order_uid = excluded.order_uid,
	transaction = excluded.transaction,
	request_id = excluded.request_id,
	currency = excluded.currency,
	provider = excluded.provider,
	amount = excluded.amount,
	payment_dt = excluded.payment_dt,
	bank = excluded.bank,
	delivery_cost = excluded.delivery_cost,
	goods_total = excluded.goods_total,
	custom_fee = excluded.custom_fee
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
    ON CONFLICT (order_uid, chrt_id) DO NOTHING
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
