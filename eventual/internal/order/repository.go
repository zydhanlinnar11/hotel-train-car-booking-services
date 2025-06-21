package order

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
)

// Repository mendefinisikan interface untuk persistensi data Order
type Repository interface {
	CreateOrder(ctx context.Context, order *Order) error
	GetOrderByID(ctx context.Context, id string) (*Order, error)
	UpdateOrder(ctx context.Context, order *Order) error
}

const (
	collectionName = "order_orders"
)

// firestoreRepository adalah implementasi konkritnya
type firestoreRepository struct {
	client *firestore.Client
}

func NewFirestoreRepository(client *firestore.Client) Repository {
	return &firestoreRepository{client: client}
}

func (r *firestoreRepository) CreateOrder(ctx context.Context, order *Order) error {
	_, err := r.client.Collection(collectionName).Doc(order.ID).Set(ctx, order)
	return err
}

func (r *firestoreRepository) GetOrderByID(ctx context.Context, id string) (*Order, error) {
	doc, err := r.client.Collection(collectionName).Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var order Order
	if err := doc.DataTo(&order); err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *firestoreRepository) UpdateOrder(ctx context.Context, order *Order) error {
	_, err := r.client.Collection(collectionName).Doc(order.ID).Set(ctx, order)
	if err != nil {
		return err
	}

	// Update updated_at
	order.UpdatedAt = time.Now()
	_, err = r.client.Collection(collectionName).Doc(order.ID).Set(ctx, order)
	return err
}
