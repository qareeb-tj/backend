package http

import (
	"context"
	"net/http"
	"qareeb/internal/domain"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandler struct {
	taxiRepo domain.TaxiServiceRepository
	timeout  time.Duration
}

func NewOrderHandler(repo domain.TaxiServiceRepository) *OrderHandler {
	return &OrderHandler{
		taxiRepo: repo,
		timeout:  10 * time.Second,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var order domain.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Генерируем ID заказа и устанавливаем время создания
	order.ID = uuid.New().String()
	order.CreatedAt = time.Now()
	order.Status = "pending"

	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), h.timeout)
	defer cancel()

	// Канал для получения первого успешного ответа
	responseChan := make(chan *domain.OrderResponse, 1)
	var wg sync.WaitGroup

	// Отправляем заказ всем службам такси
	services := h.taxiRepo.GetAllServices()
	for _, service := range services {
		wg.Add(1)
		go func(svc domain.TaxiService) {
			defer wg.Done()

			resp, err := h.taxiRepo.SendOrder(svc, order)
			if err != nil {
				return
			}

			if resp.Accepted {
				select {
				case responseChan <- resp:
				default:
					// Если канал уже заполнен, значит кто-то уже принял заказ
				}
			}
		}(service)
	}

	// Ожидаем либо первый успешный ответ, либо таймаут
	var result *domain.OrderResponse
	select {
	case result = <-responseChan:
		// Получили успешный ответ
	case <-ctx.Done():
		c.JSON(http.StatusRequestTimeout, gin.H{"error": "Timeout waiting for taxi service response"})
		return
	}

	// Отменяем все оставшиеся запросы
	cancel()

	// Ожидаем завершения всех горутин
	wg.Wait()

	if result != nil {
		c.JSON(http.StatusOK, gin.H{
			"message":  "Order accepted",
			"service":  result.ServiceName,
			"order_id": result.OrderID,
		})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "No taxi service available"})
	}
}
