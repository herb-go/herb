package notification

const (
	NotificationInstanceStatusCanceled    = -1
	NotificationInstanceStatusNew         = 0
	NotificationInstanceStatusSuccess     = 1
	NotificationInstanceStatusFail        = 2
	NotificationInstanceStatusPending     = 3
	NotificationInstanceStatusUnsupported = 10
	NotificationInstanceStatusError       = 99
)

type NotificationInstance struct {
	Notification Notification
	Sender       string
	InstanceID   string
	Recipient    string
	Status       int
	Output       string
	Logs         []string
}

func (i *NotificationInstance) NewError(err error) error {
	i.Status = NotificationInstanceStatusError
	return &NotificationError{
		Instance: i,
		Err:      err,
	}
}

func (i *NotificationInstance) IsStatusError() bool {
	return i.Status == NotificationInstanceStatusError
}
func (i *NotificationInstance) SetStatusSuccess() {
	i.Status = NotificationInstanceStatusSuccess
}
func (i *NotificationInstance) IsStatusSuccess() bool {
	return i.Status == NotificationInstanceStatusSuccess
}
func (i *NotificationInstance) SetStatusUnsupported() {
	i.Status = NotificationInstanceStatusUnsupported
}
func (i *NotificationInstance) IsStatusUnsupported() bool {
	return i.Status == NotificationInstanceStatusUnsupported
}
func NewNotificationInstance(Notification Notification) *NotificationInstance {
	return &NotificationInstance{
		Notification: Notification,
		Status:       NotificationInstanceStatusNew,
		Logs:         []string{},
	}
}
