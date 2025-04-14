package main

import (
	"log"
	"qareeb/internal/delivery/http"
	"qareeb/internal/domain"
	"qareeb/internal/infrastructure/taxi"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация репозитория такси-сервисов
	taxiServices := []domain.TaxiService{
		{Name: "Service1", Endpoint: "http://service1.example.com/order"},
		{Name: "Service2", Endpoint: "http://service2.example.com/order"},
		// Добавьте другие службы такси здесь
	}

	repo := taxi.NewRepository(taxiServices, 10*time.Second)
	handler := http.NewOrderHandler(repo)

	// Настройка маршрутизации
	r := gin.Default()
	r.POST("/order", handler.CreateOrder)

	// Запуск сервера
	log.Fatal(r.Run(":8080"))
}
