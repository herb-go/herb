package notification

type InstancesBuilder interface {
	IDGenerator() (string, error)
	SetIDGenerator(IDGenerator) error
	NewInstance(Notification) (*NotificationInstance, error)
	CloneInstance(*NotificationInstance) (*NotificationInstance, error)
	BuildInstances(Notification) ([]*NotificationInstance, error)
}

type CommonInstancesBuilder struct {
	idGenerator IDGenerator
}

func (b *CommonInstancesBuilder) IDGenerator() (string, error) {
	if b.idGenerator != nil {
		return b.idGenerator()
	}
	return DefaultIDGenerator()
}

func (b *CommonInstancesBuilder) SetIDGenerator(idGenerator IDGenerator) error {
	b.idGenerator = idGenerator
	return nil
}
func (b *CommonInstancesBuilder) NewInstance(n Notification) (*NotificationInstance, error) {
	id, err := b.IDGenerator()
	if err != nil {
		return nil, err
	}
	r, err := n.NotificationRecipient()
	if err != nil {
		return nil, err
	}
	i := NewNotificationInstance(n)
	i.Recipient = r
	i.InstanceID = id
	return i, nil
}
func (b *CommonInstancesBuilder) CloneInstance(i *NotificationInstance) (*NotificationInstance, error) {
	id, err := b.IDGenerator()
	if err != nil {
		return nil, err
	}
	clone := &NotificationInstance{
		Notification: i.Notification,
		Sender:       i.Sender,
		InstanceID:   id,
		Recipient:    i.Recipient,
		Status:       i.Status,
		Output:       i.Output,
		Logs:         make([]string, len(i.Logs)),
	}
	for k := range i.Logs {
		clone.Logs[k] = i.Logs[k]
	}
	return clone, nil

}

func (b *CommonInstancesBuilder) BuildInstances(n Notification) ([]*NotificationInstance, error) {
	ni, err := b.NewInstance(n)
	if err != nil {
		return nil, err
	}
	return []*NotificationInstance{ni}, nil
}

func NewCommonInstancesBuilder() *CommonInstancesBuilder {
	return &CommonInstancesBuilder{}
}
