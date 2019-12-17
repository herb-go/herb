package channel

import (
	"errors"
	"net/http"
	"testing"
)

func TestErrChannelUsed(t *testing.T) {
	var err error
	defer resetAll()
	channel := newTestChannel("/test")
	err = channel.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	if err != nil {
		t.Fatal(err)
	}
	err = channel.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	if err == nil || !errors.Is(err, ErrChannelUsed) {
		t.Fatal(err)
	}
}

func TestErrs(t *testing.T) {
	var err error
	defer resetAll()
	channel := newTestChannel("/test")
	err = channel.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	if err != nil {
		t.Fatal(err)
	}
	err = channel.Stop()
	if err == nil || !errors.Is(err, ErrChannelStopped) {
		t.Fatal(err)
	}
	err = channel.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer channel.Stop()
	err = channel.Start()
	if err == nil || !errors.Is(err, ErrChannelStarted) {
		t.Fatal(err)
	}
}
func TestErrChannelNotRegistered(t *testing.T) {
	var err error
	defer resetAll()
	channel := newTestChannel("/test")
	err = channel.Start()
	if err == nil || !errors.Is(err, ErrChannelNotRegistered) {
		t.Fatal(err)
	}
	err = channel.Stop()
	if err == nil || !errors.Is(err, ErrChannelNotRegistered) {
		t.Fatal(err)
	}
}

func TestHandleRecover(t *testing.T) {
	var err error
	defer resetAll()
	channel := newTestChannel("/test")
	server := GetServer(&channel.ListenerConfig)
	server.mux.Handle("/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))
	err = channel.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	if err == nil {
		t.Fatal(err)
	}
}
