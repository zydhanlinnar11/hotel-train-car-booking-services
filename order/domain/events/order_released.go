package events

import "time"

type OrderReleased struct {
	Id    string
	Cause int
	Time  time.Time
}

func NewOrderReleased(id string, cause int) *OrderReleased {
	return &OrderReleased{
		Id:    id,
		Cause: cause,
		Time:  time.Now(),
	}
}

func (o *OrderReleased) Name() string {
	return "order_released"
}

func (o *OrderReleased) Timestamp() time.Time {
	return o.Time
}
