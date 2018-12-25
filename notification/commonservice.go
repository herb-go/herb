package notification

var DefaultRecover = func() {
}

type CommonService struct {
	registeredFactories        map[string]NotificationFactory
	registeredSender           map[string][]Sender
	registeredInstancesBuilder map[string]InstancesBuilder
	idGenerator                IDGenerator
	c                          chan Notification
	closeChan                  chan bool
	recover                    func()
}

func NewCommonService() *CommonService {
	s := &CommonService{
		registeredFactories:        map[string]NotificationFactory{},
		registeredSender:           map[string][]Sender{},
		registeredInstancesBuilder: map[string]InstancesBuilder{},
		idGenerator:                DefaultIDGenerator,
		c:                          make(chan Notification),
		closeChan:                  make(chan bool),
		recover:                    DefaultRecover,
	}
	return s
}
func (m *CommonService) Recover() func() {
	return m.recover
}
func (m *CommonService) SetRecover(r func()) {
	m.recover = r
}
func (m *CommonService) Start() error {
	m.closeChan = make(chan bool)
	m.c = make(chan Notification)
	c := m.NotificationChan()
	go func() {
		for {
			select {
			case n := <-c:
				go func() {
					defer m.Recover()()
					err := SendNotificationByService(m, n)
					if err != nil {
						panic(err)
					}
				}()
			}
		}
	}()
	return nil
}

func (m *CommonService) Stop() error {
	return nil
}

func (m *CommonService) RegisterInstancesBuilder(notificationtype string, builder InstancesBuilder) error {
	m.registeredInstancesBuilder[notificationtype] = builder
	return nil
}

func (m *CommonService) InstancesBuildersByType(notificationtype string) (InstancesBuilder, error) {
	b := m.registeredInstancesBuilder[notificationtype]
	if b == nil && m.registeredInstancesBuilder[NotificationTypeDefault] != nil {
		return m.registeredInstancesBuilder[NotificationTypeDefault], nil
	}
	return b, nil
}

func (m *CommonService) NotificationChan() chan Notification {
	return m.c
}
func (m *CommonService) RegisterNotificationType(name string, f NotificationFactory) error {
	m.registeredFactories[name] = f
	return nil
}

func (m *CommonService) RegisterSender(notificationtype string, s Sender) error {
	_, ok := m.registeredSender[notificationtype]
	if ok == false {
		m.registeredSender[notificationtype] = []Sender{s}
	} else {
		m.registeredSender[notificationtype] = append(m.registeredSender[notificationtype], s)
	}
	return nil
}

func (m *CommonService) SendersByType(notificationtype string) ([]Sender, error) {
	return m.registeredSender[notificationtype], nil
}

func init() {
	s := NewCommonService()
	DefaultService = s
	DefaultService.RegisterInstancesBuilder(NotificationTypeDefault, NewCommonInstancesBuilder())
	DefaultNotifier = DefaultService
}
