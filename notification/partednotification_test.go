package notification

import "testing"

func TestPartedNotification(t *testing.T) {
	notification1, err := NewPartedNotificationWithID()
	if err != nil {
		t.Fatal(err)
	}
	nid1, err := notification1.NotificationID()
	if err != nil {
		t.Fatal(err)
	}
	if nid1 == "" {
		t.Error(nid1)
	}
	notification2, err := NewPartedNotificationWithID()
	if err != nil {
		t.Error(err)
	}
	nid2, err := notification2.NotificationID()
	if err != nil {
		t.Fatal(err)
	}
	if nid2 == "" {
		t.Error(nid2)
	}
	if nid1 == nid2 {
		t.Error(nid1)
	}
}
