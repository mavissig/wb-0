package model

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       string `json:"amount"`
	PaymentDt    string `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost string `json:"delivery_cost"`
	GoodsTotal   string `json:"goods_total"`
	CustomFee    string `json:"custom_fee"`
}

type Item struct {
	ChrtID      string `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       string `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        string `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  string `json:"total_price"`
	NmID        string `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      string `json:"status"`
}

type Order struct {
	OrderUID          string   `json:"order_uid"`
	TrackNumber       string   `json:"track_number"`
	Entry             string   `json:"entry"`
	Delivery          Delivery `json:"delivery"`
	Payment           Payment  `json:"payment"`
	Items             []Item   `json:"items"`
	Locale            string   `json:"locale"`
	InternalSignature string   `json:"internal_signature"`
	CustomerID        string   `json:"customer_id"`
	DeliveryService   string   `json:"delivery_service"`
	ShardKey          string   `json:"shardkey"`
	SMID              string   `json:"sm_id"`
	DateCreated       string   `json:"date_created"`
	OofShard          string   `json:"oof_shard"`
}
