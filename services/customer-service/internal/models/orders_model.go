package models

type OrdersDetailRequest struct {
	CustomerID int64 `json:"customer_id" param:"customer_id"`
}

type Order struct {
	OrderID    int64   `json:"order_id"`
	CustomerID int64   `json:"customer_id"`
	OrderNo    string  `json:"order_no"`
	GrandTotal float64 `json:"grand_total"`
	Status     string  `json:"status"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type OrderItem struct {
	OrderItemID int64   `json:"order_item_id"`
	OrderID     int64   `json:"order_id"`
	SKU         string  `json:"sku"`
	Qty         int     `json:"qty"`
	Amount      float64 `json:"amount"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type OrdersDetailResponse struct {
	Order
	Items []OrderItem `json:"items"`
}
