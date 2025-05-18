package events

import "time"

type Event interface {
	Name() string
	Timestamp() time.Time
}
