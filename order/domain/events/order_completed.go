package events

import "time"

type OrderCompleted struct {
	Id   string
	Time time.Time
}

func NewOrderCompleted(id string) *OrderCompleted {
	return &OrderCompleted{
		Id:   id,
		Time: time.Now(),
	}
}

func (o *OrderCompleted) Name() string {
	return "order_completed"
}

func (o *OrderCompleted) Timestamp() time.Time {
	return o.Time
}
