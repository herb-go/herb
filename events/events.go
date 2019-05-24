package events

import (
	"sync"
)

//NewEvent create new event.
func NewEvent() *Event {
	return &Event{}
}

//Type event type
type Type string

//NewEvent create event wth type.
func (t Type) NewEvent() *Event {
	return &Event{
		Type: t,
	}
}

// Hanlder event handler
type Hanlder func(*Event)

//Event event struct
type Event struct {
	Type   Type
	Target string
	Data   interface{}
}

//WithTarget set event target and return event
func (e *Event) WithTarget(target string) *Event {
	e.Target = target
	return e
}

//WithData set event data and return event
func (e *Event) WithData(data interface{}) *Event {
	e.Data = data
	return e
}

//WithType set event type and return event
func (e *Event) WithType(t Type) *Event {
	e.Type = t
	return e
}

//EventService event service which handle evetns
type EventService struct {
	Hanlders map[Type][]Hanlder
	lock     sync.RWMutex
}

//NewEvent create new event
func (e *EventService) NewEvent() *Event {
	return &Event{}
}

//Emit emit event to event service
func (e *EventService) Emit(event *Event) bool {
	e.lock.RLock()
	defer e.lock.RUnlock()
	if e.Hanlders[event.Type] == nil {
		return false
	}
	for k := range e.Hanlders[event.Type] {
		go e.Hanlders[event.Type][k](event)
	}
	return true

}

//On register event handlter to given type.
func (e *EventService) On(t Type, hanlder Hanlder) {
	e.lock.Lock()
	defer e.lock.Unlock()

	if e.Hanlders[t] == nil {
		e.Hanlders[t] = []Hanlder{}
	}
	e.Hanlders[t] = append(e.Hanlders[t], hanlder)
}

//New create New EventsService
func New() *EventService {
	e := &EventService{
		Hanlders: map[Type][]Hanlder{},
	}
	return e
}

//WrapEmit return a default service event emitter with given type
func WrapEmit(t Type) func(*Event) bool {
	return func(e *Event) bool {
		if e == nil {
			e = NewEvent()
		}
		e.Type = t
		return Emit(e)
	}
}

//WrapOn return a defalut service hanlder registeror of give event type
func WrapOn(t Type) func(Hanlder) {
	return func(hanlder Hanlder) {
		On(t, hanlder)
	}
}

//DefaultEventService default event service.
var DefaultEventService = New()

//On register default service event handlter to given type.
func On(t Type, hanlder Hanlder) {
	DefaultEventService.On(t, hanlder)
}

//Emit emit event to default event service
func Emit(event *Event) bool {
	return DefaultEventService.Emit(event)
}
