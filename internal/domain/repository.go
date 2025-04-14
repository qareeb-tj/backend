package domain

type TaxiServiceRepository interface {
	GetAllServices() []TaxiService
	SendOrder(service TaxiService, order Order) (*OrderResponse, error)
}
