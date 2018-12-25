package notificationfilter

import "github.com/herb-go/herb/notification"

func Wrap(filter Filter, sender notification.Sender) *FilterWrapper {
	return &FilterWrapper{
		Sender: sender,
		Filter: filter,
	}
}

type FilterWrapper struct {
	notification.Sender
	Filter Filter
}

func (w *FilterWrapper) SendNotification(ni *notification.NotificationInstance) error {
	return w.Filter(ni, w.Sender.SendNotification)
}

type Filter func(instance *notification.NotificationInstance, next func(*notification.NotificationInstance) error) error

func (f Filter) Wrap(sender notification.Sender) *FilterWrapper {
	return Wrap(f, sender)
}

type RecipientConvertor func(recipient string) (string, error)

var FilterRecipient = func(convertor RecipientConvertor) Filter {
	return func(instance *notification.NotificationInstance, next func(*notification.NotificationInstance) error) error {
		recipient, err := instance.Notification.NotificationRecipient()
		if err != nil {
			return err
		}
		id, err := convertor(recipient)
		if err != nil {
			return err
		}
		err = instance.Notification.SetNotificationRecipient(id)
		if err != nil {
			return err
		}
		return next(instance)
	}
}
