package service

import "github.com/Troublor/go-trash/errs"

var Events []*Event

type Event struct {
	Name      string
	listeners []func(event Event)
}

func newEvent(eventName string) (*Event, error) {
	_, err := GetEvent(eventName)
	if err == nil {
		return nil, errs.EventExistError
	}
	event := &Event{Name: eventName, listeners: make([]func(event Event), 0, 5)}
	return event, nil
}

func (event Event) Happen() {
	for _, handler := range event.listeners {
		handler(event)
	}
}

func GetEvent(eventName string) (*Event, error) {
	for _, event := range Events {
		if event.Name == eventName {
			return event, nil
		}
	}
	return nil, errs.EventNotExistError
}

func SubscribeEvent(eventName string, handler func(event Event)) error {
	event, err := GetEvent(eventName)
	if err != nil {
		return err
	}
	event.listeners = append(event.listeners, handler)
	return nil
}

func MustSubscribeEvent(eventName string, handler func(event Event)) {
	err := SubscribeEvent(eventName, handler)
	if err != nil {
		panic(err)
	}
}

func EventHappen(eventName string) error {
	event, err := GetEvent(eventName)
	if err != nil {
		return err
	}
	event.Happen()
	return nil
}

func MustEventHappen(eventName string) {
	err := EventHappen(eventName)
	if err != nil {
		panic(err)
	}
}

func CreateEvent(eventName string) error {
	event, err := newEvent(eventName)
	if err != nil {
		return err
	}
	Events = append(Events, event)
	return nil
}
