package domain

import "time"

type Order struct {
	ID        string    `json:"id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
}

type TaxiService struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
}

type OrderResponse struct {
	ServiceName string `json:"service_name"`
	OrderID     string `json:"order_id"`
	Accepted    bool   `json:"accepted"`
	Message     string `json:"message,omitempty"`
}
