package events

import "time"

type OrderBooked struct {
	Id   string
	Time time.Time
}

func NewOrderBooked(id string) *OrderBooked {
	return &OrderBooked{
		Id:   id,
		Time: time.Now(),
	}
}

func (o *OrderBooked) Name() string {
	return "order_booked"
}

func (o *OrderBooked) Timestamp() time.Time {
	return o.Time
}
