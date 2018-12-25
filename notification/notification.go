package notification

import "sync"

type Notification interface {
	NotificationID() (string, error)
	SetNotificationID(string) error
	NotificationType() (string, error)
	SetNotificationType(string) error
	NotificationRecipient() (string, error)
	SetNotificationRecipient(string) error
	NotificationAuthor() (string, error)
	SetNotificationAuthor(string) error
	LockNotification()
	UnlockNotification()
	RLockNotification()
	RUnlockNotification()
}

type CommonNotification struct {
	ID        string
	Type      string
	Recipient string
	Author    string
	lock      sync.RWMutex
}

func (m *CommonNotification) LockNotification() {
	m.lock.Lock()
}
func (m *CommonNotification) UnlockNotification() {
	m.lock.Unlock()
}
func (m *CommonNotification) RLockNotification() {
	m.lock.RLock()
}
func (m *CommonNotification) RUnlockNotification() {
	m.lock.RUnlock()
}
func (m *CommonNotification) NotificationAuthor() (string, error) {
	return m.Author, nil
}

func (m *CommonNotification) SetNotificationAuthor(author string) error {
	m.Author = author
	return nil
}

func (m *CommonNotification) NotificationID() (string, error) {
	return m.ID, nil
}

func (m *CommonNotification) SetNotificationID(id string) error {
	m.ID = id
	return nil
}

func (m *CommonNotification) NotificationType() (string, error) {
	return m.Type, nil
}

func (m *CommonNotification) SetNotificationType(notificationtype string) error {
	m.Type = notificationtype
	return nil
}

func (m *CommonNotification) NotificationRecipient() (string, error) {
	return m.Recipient, nil
}

func (m *CommonNotification) SetNotificationRecipient(recipient string) error {
	m.Recipient = recipient
	return nil
}
