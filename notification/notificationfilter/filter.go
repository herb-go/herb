package notificationfilter

import "github.com/herb-go/herb/notification"

func Wrap(filter Filter, sender notification.Sender) *WrappedFilter {
	return &WrappedFilter{
		Sender: sender,
		Filter: filter,
	}
}

type WrappedFilter struct {
	notification.Sender
	Filter Filter
}

func (w *WrappedFilter) SendNotification(ni *notification.NotificationInstance) error {
	return w.Filter(ni, w.Sender.SendNotification)
}

type Filter func(instance *notification.NotificationInstance, next func(*notification.NotificationInstance) error) error

func (f Filter) Wrap(sender notification.Sender) *WrappedFilter {
	return Wrap(f, sender)
}

type RecipientConvertor func(recipient string) (string, error)

var RecipientConvertorWrapper = func(convertor RecipientConvertor) Filter {
	return func(instance *notification.NotificationInstance, next func(*notification.NotificationInstance) error) error {
		recipient, err := instance.Notification.NotificationRecipient()
		if err != nil {
			return err
		}
		id, err := convertor(recipient)
		if err != nil {
			return err
		}
		if id == "" {
			instance.SetStatusCanceled()
			return nil
		}
		err = instance.Notification.SetNotificationRecipient(id)
		if err != nil {
			return err
		}
		return next(instance)
	}
}
