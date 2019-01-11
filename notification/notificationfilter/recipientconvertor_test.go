package notificationfilter

import "testing"
import "github.com/herb-go/herb/notification"

type userMaps map[string]string

func (m *userMaps) RecipientConvertor(recipient string) (string, error) {
	if m == nil {
		return "", nil
	}
	return (*m)[recipient], nil
}

var UserMapsTest = userMaps{
	"raw1": "test1",
	"raw2": "test2",
}

type ChanSender struct {
	ID string
	C  chan *notification.NotificationInstance
}

func (c *ChanSender) SendNotification(n *notification.NotificationInstance) error {
	c.C <- n
	return nil
}
func (c *ChanSender) Name() string {
	return c.ID
}

func newChanSender() *ChanSender {
	id, err := notification.DefaultIDGenerator()
	if err != nil {
		panic(err)
	}
	return &ChanSender{
		ID: id,
		C:  make(chan *notification.NotificationInstance, 10),
	}
}

func Test_FilterRecipient(t *testing.T) {

}
