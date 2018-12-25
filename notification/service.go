package notification

type NotificationFactory func() (Notification, error)

const NotificationTypeDefault = ""

type Service interface {
	RegisterNotificationType(name string, f NotificationFactory) error
	RegisterSender(notificationtype string, s Sender) error
	SendersByType(notificationtype string) ([]Sender, error)
	Recover() func()
	SetRecover(func())
	RegisterInstancesBuilder(notificationtype string, builder InstancesBuilder) error
	InstancesBuildersByType(notificationtype string) (InstancesBuilder, error)
	Start() error
	Stop() error
	Notifier
}

var SendNotificationByService = func(m Service, n Notification) error {
	nt, err := n.NotificationType()
	if err != nil {
		return err
	}
	builder, err := m.InstancesBuildersByType(nt)
	if err != nil {
		return err
	}
	notificationInstances, err := builder.BuildInstances(n)
	if err != nil {
		return err
	}
	Senders, err := m.SendersByType(nt)
	if err != nil {
		return err
	}
	for i := range notificationInstances {
		for k := range Senders {
			ni, err := builder.CloneInstance(notificationInstances[i])
			if err != nil {
				return err
			}
			go func(sender Sender) {
				defer m.Recover()()
				ni.Sender = sender.Name()
				err = sender.SendNotification(ni)
				if err != nil {
					panic(err)
				}
			}(Senders[k])
		}
	}
	return nil
}

var DefaultService Service
