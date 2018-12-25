package notification

import (
	"errors"
	"strings"
	"testing"
)

func TestIntanceError(t *testing.T) {
	defer func() {
		r := recover()
		err, ok := r.(error)
		if ok == false {
			t.Fatal(ok)
		}
		nerr, ok := r.(*NotificationError)
		if ok == false {
			t.Fatal(ok)
		}
		if !nerr.Instance.IsStatusError() {
			t.Fatal(nerr.Instance)
		}
		output := err.Error()
		if !strings.Contains(output, "test error") {
			t.Fatal(output)
		}
		if !strings.Contains(output, "testid") {
			t.Fatal(output)
		}
		if !strings.Contains(output, "testsender") {
			t.Fatal(output)
		}
		if !strings.Contains(output, "testlog1") {
			t.Fatal(output)
		}
		if !strings.Contains(output, "testlog2") {
			t.Fatal(output)
		}
	}()
	err := errors.New("test error")
	n := NewNotificationInstance(NewPartedNotification())
	n.InstanceID = "testid"
	n.Sender = "testsender"
	n.Logs = []string{"testlog1", "testlog2"}
	panic(n.NewError(err))
}
