package model

import "time"

type Delivery struct {
	Name    string `json:"name" db:"name"`
	Phone   string `json:"phone" db:"phone"`
	Zip     string `json:"zip" db:"zip"`
	City    string `json:"city" db:"city"`
	Address string `json:"address" db:"address"`
	Region  string `json:"region" db:"region"`
	Email   string `json:"email" db:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction" db:"transaction"`
	RequestID    string `json:"request_id" db:"request_id"`
	Currency     string `json:"currency" db:"currency"`
	Provider     string `json:"provider" db:"provider"`
	Amount       int    `json:"amount" db:"amount"`
	PaymentDt    int    `json:"payment_dt" db:"payment_dt"`
	Bank         string `json:"bank" db:"bank"`
	DeliveryCost int    `json:"delivery_cost" db:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total" db:"goods_total"`
	CustomFee    int    `json:"custom_fee" db:"custom_fee"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id" db:"chrt_id"`
	TrackNumber string `json:"track_number" db:"track_number"`
	Price       int    `json:"price" db:"price"`
	Rid         string `json:"rid" db:"rid"`
	Name        string `json:"name" db:"name"`
	Sale        int    `json:"sale" db:"sale"`
	Size        string `json:"size" db:"size"`
	TotalPrice  int    `json:"total_price" db:"total_price"`
	NmID        int    `json:"nm_id" db:"nm_id"`
	Brand       string `json:"-"`
	Status      int    `json:"-"`
}

type Order struct {
	OrderUID          string    `db:"order_uid" json:"order_uid"`
	TrackNumber       string    `db:"track_number" json:"track_number"`
	Entry             string    `db:"entry" json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `db:"locale" json:"locale"`
	InternalSignature string    `db:"internal_signature" json:"internal_signature"`
	CustomerID        string    `db:"customer_id" json:"customer_id"`
	DeliveryService   string    `db:"delivery_service" json:"delivery_service"`
	ShardKey          string    `db:"shardkey" json:"shardkey"`
	SMID              int       `db:"sm_id" json:"sm_id"`
	DateCreated       time.Time `db:"date_created" json:"date_created"`
	OofShard          string    `db:"oof_shard" json:"oof_shard"`
}
