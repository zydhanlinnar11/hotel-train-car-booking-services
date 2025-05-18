package infrastructures

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/zydhanlinnar11/sister/ec/order/domain/entities"
	"github.com/zydhanlinnar11/sister/ec/order/domain/events"
	"github.com/zydhanlinnar11/sister/ec/order/domain/repositories"
	"github.com/zydhanlinnar11/sister/ec/order/domain/values"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FirestoreOrderRepository struct {
	client           *firestore.Client
	orderIdValidator values.OrderIdValidator
}

func NewFirestoreOrderRepository(client *firestore.Client, orderIdValidator values.OrderIdValidator) repositories.Order {
	return &FirestoreOrderRepository{
		client:           client,
		orderIdValidator: orderIdValidator,
	}
}

// FindById implements repositories.Order.
func (r *FirestoreOrderRepository) FindById(ctx context.Context, id string) (*entities.Order, error) {
	doc, err := r.client.Collection(OrderCollection).Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}

	var data Order
	if err := doc.DataTo(&data); err != nil {
		return nil, err
	}

	vi, err := values.NewOrderId(data.Id, r.orderIdValidator)
	if err != nil {
		return nil, err
	}

	hotelRoomStartDate, err := values.NewDate(data.HotelRoomStartDate)
	if err != nil {
		return nil, err
	}
	hotelRoomEndDate, err := values.NewDate(data.HotelRoomEndDate)
	if err != nil {
		return nil, err
	}

	carStartDate, err := values.NewDate(data.CarStartDate)
	if err != nil {
		return nil, err
	}
	carEndDate, err := values.NewDate(data.CarEndDate)
	if err != nil {
		return nil, err
	}

	order := entities.NewOrder(
		vi,
		data.HotelRoomId,
		hotelRoomStartDate,
		hotelRoomEndDate,
		data.CarId,
		carStartDate,
		carEndDate,
		data.TrainSeatId,
		entities.OrderFailureCause(data.FailureCause),
		data.Version,
		entities.OrderStatus(data.Status),
	)

	return order, nil
}

func (r *FirestoreOrderRepository) convertToSerializableEvents(ctx context.Context, evs []events.Event) []Event {
	serializableEvents := make([]Event, 0)
	for _, event := range evs {
		newEvent := Event{
			Name:    event.Name(),
			Payload: event.Timestamp(),
		}
		switch e := event.(type) {
		case *events.OrderPlaced:
			newEvent.Payload = OrderPlacedEventPayload{
				OrderId:            e.Id,
				HotelRoomId:        e.HotelRoomId,
				HotelRoomStartDate: e.HotelRoomStartDate,
				HotelRoomEndDate:   e.HotelRoomEndDate,
				CarId:              e.CarId,
				CarStartDate:       e.CarStartDate,
				CarEndDate:         e.CarEndDate,
				TrainSeatId:        e.TrainSeatId,
			}
		case *events.OrderBooked:
			newEvent.Payload = OrderBookedEventPayload{
				OrderId: e.Id,
			}
		case *events.OrderCompleted:
			newEvent.Payload = OrderCompletedEventPayload{
				OrderId: e.Id,
			}
		case *events.OrderReleased:
			newEvent.Payload = OrderReleasedEventPayload{
				OrderId: e.Id,
				Cause:   e.Cause,
			}
		case *events.OrderFailed:
			newEvent.Payload = OrderFailedEventPayload{
				OrderId: e.Id,
				Cause:   e.Cause,
			}
		default:
			continue
		}

		serializableEvents = append(serializableEvents, newEvent)
	}

	return serializableEvents
}

func (r *FirestoreOrderRepository) serializeOrder(ctx context.Context, order *entities.Order) Order {
	return Order{
		Id:                 order.Id().String(),
		HotelRoomId:        order.HotelRoomId(),
		HotelRoomStartDate: order.HotelRoomStartDate().String(),
		HotelRoomEndDate:   order.HotelRoomEndDate().String(),
		CarId:              order.CarId(),
		CarStartDate:       order.CarStartDate().String(),
		CarEndDate:         order.CarEndDate().String(),
		TrainSeatId:        order.TrainSeatId(),
		FailureCause:       order.FailureCause().Int(),
		Version:            order.Version(),
		Status:             int(order.Status()),
	}
}

func (r *FirestoreOrderRepository) Save(ctx context.Context, order *entities.Order) error {
	events := r.convertToSerializableEvents(ctx, order.Events())
	serializedOrder := r.serializeOrder(ctx, order)
	orderRef := r.client.Collection(OrderCollection).Doc(serializedOrder.Id)

	err := r.client.RunTransaction(ctx, func(ctx context.Context, t *firestore.Transaction) error {
		orderDoc, err := t.Get(orderRef)
		if err != nil && status.Code(err) != codes.NotFound {
			return err
		}

		if orderDoc != nil {
			var existingOrder Order
			if err := orderDoc.DataTo(&existingOrder); err != nil {
				return err
			}

			if existingOrder.Version >= serializedOrder.Version {
				return repositories.ErrOrderVersionMismatch
			}

			serializedOrder.Version = existingOrder.Version + 1
		}

		if err := t.Set(orderRef, serializedOrder); err != nil {
			return err
		}

		for _, event := range events {
			eventRef := r.client.Collection(EventCollection).NewDoc()

			if err := t.Set(eventRef, event); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
