package main

import (
	"context"
	calendar "github.com/anton.okolelov/otus-calendar-grpc/internal/grpc"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	address = "localhost:50051"
)

func main() {

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := calendar.NewCalendarClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	start, err := ptypes.TimestampProto(time.Now().Add(time.Hour))
	end, err := ptypes.TimestampProto(time.Now().Add(2 * time.Hour))

	if err != nil {
		log.Printf("%v", err)
	}

	calendarEvent := &calendar.Event{
		Start:   start,
		End:     end,
		Payload: "test",
	}

	eventId, err := c.AddEvent(ctx, calendarEvent)

	if err != nil {
		log.Printf("Error adding event: %v", err)
	}

	log.Printf("Event added %v", eventId.Id)

	calendarEvent.Payload = "test2"

	eventId, err = c.UpdateEvent(ctx, &calendar.EventUpdateInfo{
		EventId: eventId.Id,
		Event:   calendarEvent,
	})

	if err != nil {
		log.Printf("Error adding event: %v", err)
	}

	log.Printf("Event updated %v", eventId.Id)

	events, err := c.GetEvents(ctx, &calendar.GetEventsRequest{})
	for _, event := range events.Events {
		log.Printf("Event: %v", event)
	}

	eventId, err = c.DeleteEvent(ctx, eventId)

	if err != nil {
		log.Printf("Error adding event: %v", err)
	}

	log.Printf("Event deleted: %v", eventId.Id)
}
