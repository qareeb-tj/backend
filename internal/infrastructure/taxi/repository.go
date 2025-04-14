package taxi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"qareeb/internal/domain"
	"time"
)

type Repository struct {
	services []domain.TaxiService
	client   *http.Client
	timeout  time.Duration
}

func NewRepository(services []domain.TaxiService, timeout time.Duration) *Repository {
	return &Repository{
		services: services,
		client:   &http.Client{Timeout: timeout},
		timeout:  timeout,
	}
}

func (r *Repository) GetAllServices() []domain.TaxiService {
	return r.services
}

func (r *Repository) SendOrder(service domain.TaxiService, order domain.Order) (*domain.OrderResponse, error) {
	jsonData, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", service.Endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var orderResp domain.OrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&orderResp); err != nil {
		return nil, err
	}

	return &orderResp, nil
}
