package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/zydhanlinnar11/sister/ec/order/app/usecases"
	"github.com/zydhanlinnar11/sister/ec/order/infrastructures"
)

const (
	defaultPort = "8080"
)

type PlaceOrderRequest struct {
	HotelRoomId        string `json:"hotel_room_id"`
	HotelRoomStartDate string `json:"hotel_room_start_date"`
	HotelRoomEndDate   string `json:"hotel_room_end_date"`
	CarId              string `json:"car_id"`
	CarStartDate       string `json:"car_start_date"`
	CarEndDate         string `json:"car_end_date"`
	TrainSeatId        string `json:"train_seat_id"`
}

// Example request body
// {
// 	"hotel_room_id": "123",
// 	"hotel_room_start_date": "01-01-2025",
// 	"hotel_room_end_date": "02-01-2025",
// 	"car_id": "456",
// 	"car_start_date": "01-01-2025",
// 	"car_end_date": "02-01-2025",
// 	"train_seat_id": "789",
// }

func main() {
	godotenv.Load()
	firestoreClient, err := firestore.NewClient(context.Background(), os.Getenv("GOOGLE_PROJECT_ID"))
	if err != nil {
		log.Fatalf("failed to create firestore client: %v", err)
	}

	orderIdValidatorGenerator := infrastructures.NewOrderIdValidatorGenerator()
	orderRepository := infrastructures.NewFirestoreOrderRepository(firestoreClient, orderIdValidatorGenerator)

	placeOrderUc := &usecases.PlaceOrderUseCase{
		OrderRepo:  orderRepository,
		OrderIdGen: orderIdValidatorGenerator,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/orders/place", func(w http.ResponseWriter, r *http.Request) {
		var req PlaceOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err := placeOrderUc.Execute(r.Context(), usecases.PlaceOrderRequest{
			HotelRoomId:        req.HotelRoomId,
			HotelRoomStartDate: req.HotelRoomStartDate,
			HotelRoomEndDate:   req.HotelRoomEndDate,
			CarId:              req.CarId,
			CarStartDate:       req.CarStartDate,
			CarEndDate:         req.CarEndDate,
			TrainSeatId:        req.TrainSeatId,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"order_id": res.OrderId,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Printf("server is running on port %s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}
