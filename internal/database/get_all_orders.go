package database

import (
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jmoiron/sqlx"
	"wb/internal/model"
)

func (database *Database) GetAllOrders() (map[string]model.Order, error) {
	fmt.Println("GetAllOrders")

	type Temp struct {
		model.Order
		model.Delivery
		model.Payment
		model.Item
	}

	var temps []Temp
	query := `
    SELECT 
        Orders.order_uid,
        Orders.track_number,
        Orders.entry,
        Orders.locale,
        Orders.internal_signature,
        Orders.customer_id,
        Orders.delivery_service,
        Orders.shardkey,
        Orders.sm_id,
        Orders.date_created,
        Orders.oof_shard,

        Delivery.name,
        Delivery.phone,
        Delivery.zip,
        Delivery.city,
        Delivery.address,
        Delivery.region,
        Delivery.email,

        Payment.transaction,
        Payment.request_id,
        Payment.currency,
        Payment.provider,
        Payment.amount,
        Payment.payment_dt,
        Payment.bank,
        Payment.delivery_cost,
        Payment.goods_total,
        Payment.custom_fee,

        Items.chrt_id,
        Items.track_number,
        Items.price,
        Items.rid,
        Items.name,
        Items.sale,
        Items.size,
        Items.total_price,
        Items.nm_id,
        Items.brand,
        Items.status

    FROM Orders
    LEFT JOIN Delivery ON Orders.order_uid = Delivery.order_uid
    LEFT JOIN Payment ON Orders.order_uid = Payment.order_uid
    LEFT JOIN OrderItems ON Orders.order_uid = OrderItems.order_uid
    LEFT JOIN Items ON OrderItems.chrt_id = Items.chrt_id;
`

	err := database.db.Select(&temps, query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	orderMap := make(map[string]model.Order)
	for _, temp := range temps {
		order, exists := orderMap[temp.OrderUID]
		if !exists {
			order = temp.Order
			order.Delivery = temp.Delivery
			order.Payment = temp.Payment
			order.Items = []model.Item{}
			orderMap[order.OrderUID] = order
		}
		order.Items = append(order.Items, temp.Item)
	}

	fmt.Println("[init OrderUID]: ", orderMap["testid"].OrderUID)
	return orderMap, nil
}
