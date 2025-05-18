package entities

import (
	"github.com/zydhanlinnar11/sister/ec/order/domain/events"
	"github.com/zydhanlinnar11/sister/ec/order/domain/values"
)

type OrderStatus int

const (
	// placed, booked, released, completed, failed
	OrderStatusPlaced OrderStatus = iota + 1
	OrderStatusBooked
	OrderStatusReleased
	OrderStatusCompleted
	OrderStatusFailed
)

type OrderReleaseCause int

func NewOrderReleaseCause(cause int) OrderReleaseCause {
	switch cause {
	case OrderReleaseCauseHotelRoomNotAvailable.Int():
		return OrderReleaseCauseHotelRoomNotAvailable
	case OrderReleaseCauseCarNotAvailable.Int():
		return OrderReleaseCauseCarNotAvailable
	case OrderReleaseCauseTrainSeatNotAvailable.Int():
		return OrderReleaseCauseTrainSeatNotAvailable
	default:
		return OrderReleaseCauseOther
	}
}

func (c OrderReleaseCause) Int() int {
	return int(c)
}

const (
	OrderReleaseCauseHotelRoomNotAvailable OrderReleaseCause = iota + 1
	OrderReleaseCauseCarNotAvailable
	OrderReleaseCauseTrainSeatNotAvailable
	OrderReleaseCauseOther
)

type OrderFailureCause = OrderReleaseCause

type Order struct {
	id values.OrderId

	hotelRoomId        string
	hotelRoomStartDate values.Date
	hotelRoomEndDate   values.Date

	carId        string
	carStartDate values.Date
	carEndDate   values.Date

	trainSeatId string

	failureCause OrderFailureCause

	version int
	status  OrderStatus
	events  []events.Event
}

func NewOrder(
	id values.OrderId,

	hotelRoomId string,
	hotelRoomStartDate values.Date,
	hotelRoomEndDate values.Date,

	carId string,
	carStartDate values.Date,
	carEndDate values.Date,

	trainSeatId string,

	failureCause OrderFailureCause,

	version int,
	status OrderStatus,
) *Order {
	return &Order{
		id: id,

		hotelRoomId:        hotelRoomId,
		hotelRoomStartDate: hotelRoomStartDate,
		hotelRoomEndDate:   hotelRoomEndDate,

		carId:        carId,
		carStartDate: carStartDate,
		carEndDate:   carEndDate,

		trainSeatId: trainSeatId,

		failureCause: failureCause,

		version: version,
		status:  status,
		events:  make([]events.Event, 0),
	}
}

func PlaceOrder(
	idGen values.OrderIdGenerator,
	hotelRoomId string,
	hotelRoomStartDate values.Date,
	hotelRoomEndDate values.Date,
	carId string,
	carStartDate values.Date,
	carEndDate values.Date,
	trainSeatId string,
) *Order {
	id := idGen.Generate()

	o := NewOrder(
		id,
		hotelRoomId,
		hotelRoomStartDate,
		hotelRoomEndDate,
		carId,
		carStartDate,
		carEndDate,
		trainSeatId,
		OrderReleaseCauseOther,
		0,
		OrderStatusPlaced,
	)

	return o
}

func (o *Order) Id() values.OrderId {
	return o.id
}

func (o *Order) HotelRoomId() string {
	return o.hotelRoomId
}

func (o *Order) HotelRoomStartDate() values.Date {
	return o.hotelRoomStartDate
}

func (o *Order) HotelRoomEndDate() values.Date {
	return o.hotelRoomEndDate
}

func (o *Order) CarId() string {
	return o.carId
}

func (o *Order) CarStartDate() values.Date {
	return o.carStartDate
}

func (o *Order) CarEndDate() values.Date {
	return o.carEndDate
}

func (o *Order) TrainSeatId() string {
	return o.trainSeatId
}

func (o *Order) Version() int {
	return o.version
}

func (o *Order) Status() OrderStatus {
	return o.status
}

func (o *Order) Events() []events.Event {
	return o.events
}

func (o *Order) Release(cause OrderReleaseCause) {
	o.status = OrderStatusReleased
	o.events = append(o.events, events.NewOrderReleased(o.id.String(), cause.Int()))
}

func (o *Order) Fail(cause OrderFailureCause) {
	o.status = OrderStatusFailed
	o.failureCause = cause
	o.events = append(o.events, events.NewOrderFailed(o.id.String(), cause.Int()))
}

func (o *Order) Book() {
	o.status = OrderStatusBooked
	o.events = append(o.events, events.NewOrderBooked(o.id.String()))
}

func (o *Order) Complete() {
	o.status = OrderStatusCompleted
	o.events = append(o.events, events.NewOrderCompleted(o.id.String()))
}
