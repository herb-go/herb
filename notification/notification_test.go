package notification

import (
	"testing"
	"time"
)

type ChanSender struct {
	ID string
	C  chan *NotificationInstance
}

func (c *ChanSender) SendNotification(n *NotificationInstance) error {
	c.C <- n
	return nil
}
func (c *ChanSender) Name() string {
	return c.ID
}

func newChanSender() *ChanSender {
	id, err := DefaultIDGenerator()
	if err != nil {
		panic(err)
	}
	return &ChanSender{
		ID: id,
		C:  make(chan *NotificationInstance, 10),
	}
}

func TestSend(t *testing.T) {
	n, err := NewPartedNotificationWithID()
	if err != nil {
		t.Fatal(err)
	}
	err = n.SetNotificationAuthor("AuthorTest")
	if err != nil {
		t.Fatal(err)
	}
	err = n.SetNotificationType("testtype")
	if err != nil {
		t.Fatal(err)
	}
	err = n.SetNotificationRecipient("RecipientTest")
	if err != nil {
		t.Fatal(err)
	}
	unusedSender := newChanSender()
	sender := newChanSender()
	err = DefaultService.RegisterSender("testtype", sender)
	if err != nil {
		t.Fatal(err)
	}
	err = DefaultService.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer DefaultService.Stop()
	if len(sender.C) != 0 {
		t.Error(sender.C)
	}
	if len(unusedSender.C) != 0 {
		t.Error(sender.C)
	}
	Notify(n)
	time.Sleep(time.Second)
	if len(sender.C) != 1 {
		t.Error(sender.C)
	}
	if len(unusedSender.C) != 0 {
		t.Error(sender.C)
	}
	ni := <-sender.C
	niNotificationID, err := ni.Notification.NotificationID()
	if err != nil {
		t.Fatal(err)
	}
	niNotificationAuthor, err := ni.Notification.NotificationAuthor()
	if err != nil {
		t.Fatal(err)
	}
	niNotificationRecipient, err := ni.Notification.NotificationRecipient()
	if err != nil {
		t.Fatal(err)
	}
	niNotificationType, err := ni.Notification.NotificationType()
	if err != nil {
		t.Fatal(err)
	}
	nNotificationID, err := n.NotificationID()
	if err != nil {
		t.Fatal(err)
	}
	nNotificationAuthor, err := n.NotificationAuthor()
	if err != nil {
		t.Fatal(err)
	}
	nNotificationRecipient, err := n.NotificationRecipient()
	if err != nil {
		t.Fatal(err)
	}
	nNotificationType, err := n.NotificationType()
	if err != nil {
		t.Fatal(err)
	}
	if niNotificationID != nNotificationID ||
		niNotificationAuthor != nNotificationAuthor ||
		niNotificationRecipient != nNotificationRecipient ||
		niNotificationType != nNotificationType {
		t.Error(*ni)
	}
}

func TestMutliSender(t *testing.T) {
	n, err := NewPartedNotificationWithID()
	if err != nil {
		t.Fatal(err)
	}
	err = n.SetNotificationAuthor("AuthorTest")
	if err != nil {
		t.Fatal(err)
	}
	err = n.SetNotificationType("testtype")
	if err != nil {
		t.Fatal(err)
	}
	err = n.SetNotificationRecipient("RecipientTest")
	if err != nil {
		t.Fatal(err)
	}
	unusedSender := newChanSender()
	sender := newChanSender()
	sender2 := newChanSender()
	service := NewCommonService()
	service.RegisterInstancesBuilder(NotificationTypeDefault, NewCommonInstancesBuilder())
	err = service.RegisterSender("testtype", sender)
	if err != nil {
		t.Fatal(err)
	}
	err = service.RegisterSender("testtype", sender2)
	if err != nil {
		t.Fatal(err)
	}
	err = service.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer service.Stop()
	if len(sender.C) != 0 {
		t.Error(sender.C)
	}
	if len(sender2.C) != 0 {
		t.Error(sender2.C)
	}
	if len(unusedSender.C) != 0 {
		t.Error(sender.C)
	}
	NotifyTo(service, n)
	time.Sleep(time.Second)
	if len(sender.C) != 1 {
		t.Error(sender.C)
	}
	if len(sender2.C) != 1 {
		t.Error(sender2.C)
	}
	if len(unusedSender.C) != 0 {
		t.Error(sender.C)
	}
	ni := <-sender.C
	ni2 := <-sender2.C

	niNotificationAuthor, err := ni.Notification.NotificationAuthor()
	if err != nil {
		t.Fatal(err)
	}
	niNotificationRecipient, err := ni.Notification.NotificationRecipient()
	if err != nil {
		t.Fatal(err)
	}
	niNotificationType, err := ni2.Notification.NotificationType()
	if err != nil {
		t.Fatal(err)
	}

	ni2NotificationAuthor, err := ni2.Notification.NotificationAuthor()
	if err != nil {
		t.Fatal(err)
	}
	ni2NotificationRecipient, err := ni2.Notification.NotificationRecipient()
	if err != nil {
		t.Fatal(err)
	}
	ni2NotificationType, err := ni2.Notification.NotificationType()
	if err != nil {
		t.Fatal(err)
	}
	if ni.InstanceID == ni2.InstanceID {
		t.Error(ni)
	}
	if niNotificationAuthor != ni2NotificationAuthor {
		t.Error(niNotificationAuthor, ni2NotificationAuthor)
	}
	if niNotificationRecipient != ni2NotificationRecipient {
		t.Error(niNotificationRecipient, ni2NotificationRecipient)
	}
	if niNotificationType != ni2NotificationType {
		t.Error(niNotificationType, ni2NotificationType)
	}
	if ni.Output != ni2.Output {
		t.Error(ni.Output, ni2.Output)
	}
	if ni.Sender == ni2.Sender {
		t.Error(ni.Sender, ni2.Sender)
	}
	if ni.Status != ni2.Status {
		t.Error(ni.Output, ni2.Output)
	}
	for k := range ni.Logs {
		if ni.Logs[k] != ni2.Logs[k] {
			t.Error(ni.Logs[k], ni2.Logs[k])
		}
	}
}
