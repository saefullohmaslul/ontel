package models

type DetailProfileRequest struct {
	CustomerID int `json:"customer_id" param:"customer_id"`
}

type DetailProfileResponse struct {
	CustomerID int `json:"customer_id"`
}

type CustomerRequest struct {
	CustomerID int
}

type Customer struct {
	CustomerID int
	Name       string
	Email      string
	Status     int
}
