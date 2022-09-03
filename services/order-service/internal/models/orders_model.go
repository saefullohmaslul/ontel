package models

type OrdersRequest struct {
	CustomerID int64
}

type Order struct {
	OrderID    int64
	CustomerID int64
	OrderNo    string
	GrandTotal float64
	Status     string
	CreatedAt  string
	UpdatedAt  string
}

type OrderItemsRequest struct {
	OrderID int64
}

type OrderItem struct {
	OrderItemID int
	OrderID     int
	SKU         string
	Qty         int
	Amount      float64
	CreatedAt   string
	UpdatedAt   string
}

type OrdersDetailRequest struct {
	CustomerID int64 `json:"customer_id" param:"customer_id"`
}

type OrdersDetailResponse struct {
	Order
	Items []OrderItem `json:"items"`
}
