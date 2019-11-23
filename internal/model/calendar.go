package model

import (
	"errors"
	"sync/atomic"
	"time"
)

type EventId uint32

type Event struct {
	Start   time.Time
	End     time.Time
	Payload string
}

type Calendar struct {
	counter uint32
	events  map[EventId]Event
}

func New() *Calendar {
	return &Calendar{
		counter: 0,
		events:  map[EventId]Event{},
	}
}

func (calendar *Calendar) AddEvent(event Event) EventId {
	newEventId := EventId(atomic.AddUint32(&calendar.counter, 1))
	calendar.events[newEventId] = event
	return newEventId
}

func (calendar *Calendar) UpdateEvent(id EventId, event Event) {
	calendar.events[id] = event
}

func (calendar *Calendar) DeleteEvent(id EventId) {
	delete(calendar.events, id)
}

func (calendar *Calendar) Event(id EventId) (Event, error) {
	if event, ok := calendar.events[id]; ok {
		return event, nil
	} else {
		return Event{}, errors.New("not found")
	}
}

func (calendar *Calendar) Events() map[EventId]Event {
	return calendar.events
}
