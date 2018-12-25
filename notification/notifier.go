package notification

type Notifier interface {
	NotificationChan() chan Notification
}

var DefaultNotifier Notifier

var NotifyTo = func(notifier Notifier, n Notification) {
	go func() {
		notifier.NotificationChan() <- n
	}()
}

var Notify = func(n Notification) {
	NotifyTo(DefaultNotifier, n)
}
