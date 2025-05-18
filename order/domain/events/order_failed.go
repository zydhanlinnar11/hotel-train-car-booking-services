package events

import "time"

type OrderFailed struct {
	Id    string
	Cause int
	Time  time.Time
}

func NewOrderFailed(id string, cause int) *OrderFailed {
	return &OrderFailed{
		Id:    id,
		Cause: cause,
		Time:  time.Now(),
	}
}

func (o *OrderFailed) Name() string {
	return "order_failed"
}

func (o *OrderFailed) Timestamp() time.Time {
	return o.Time
}
