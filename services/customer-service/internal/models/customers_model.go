package models

type DetailProfileRequest struct {
	CustomerID int64 `json:"customer_id" param:"customer_id"`
}

type DetailProfileResponse struct {
	Customer
	Orders []OrdersDetailResponse `json:"orders"`
}

type CustomerRequest struct {
	CustomerID int64
}

type Customer struct {
	CustomerID int64  `json:"customer_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Status     int    `json:"status"`
}
