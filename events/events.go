package events

import (
	"sync"

	"github.com/herb-go/herb/events"
)

type Type string

func (t Type) NewEvent() *Event {
	return &Event{
		Type: t,
	}
}

type Hanlder func(*Event)

type Event struct {
	Type   Type
	Target string
	Data   interface{}
}

func (e *Event) WithTarget(target string) *Event {
	e.Target = target
	return e
}

func (e *Event) WithData(data interface{}) *Event {
	e.Data = data
	return e
}

func (e *Event) WithType(t Type) *Event {
	e.Type = t
	return e
}

type Events struct {
	Hanlders  map[Type][]Hanlder
	c         chan *Event
	lock      sync.RWMutex
	listening bool
}

func (e *Events) NewEvent() *Event {
	return &Event{}
}
func (e *Events) Emit(event *Event) {
	e.c <- event
}
func (e *Events) On(t Type, hanlder Hanlder) {
	e.lock.Lock()
	defer e.lock.Unlock()
	if e.listening == false {
		go e.Listen()
	}
	if e.Hanlders[t] == nil {
		e.Hanlders[t] = []Hanlder{}
	}
	e.Hanlders[t] = append(e.Hanlders[t], hanlder)
}

func (e *Events) Trigger(event *Event) bool {
	e.lock.RLock()
	defer e.lock.RUnlock()
	if e.Hanlders[event.Type] == nil {
		return true
	}
	for k := range e.Hanlders[event.Type] {
		go e.Hanlders[event.Type][k](event)
	}
	return false
}

func (e *Events) Listen() {
	e.listening = true
	for {
		select {
		case event := <-e.c:
			e.Trigger(event)
		}
	}
}
func New() *Events {
	e := &Events{
		Hanlders: map[Type][]Hanlder{},
		c:        make(chan *Event),
	}
	return e
}

func WrapEmit(t events.Type) func(*events.Event) {
	return func(e *events.Event) {
		e.Type = t
		appevents.Emit(e)
	}
}

func WrapOn(t events.Type) func(events.Hanlder) {
	return func(hanlder events.Hanlder) {
		appevents.On(t, hanlder)
	}
}

var DefaultEvents = New()

func On(t Type, hanlder Hanlder) {
	DefaultEvents.On(t, hanlder)
}

func Emit(event *Event) {
	DefaultEvents.Emit(event)
}
func init() {
	go DefaultEvents.Listen()
}
