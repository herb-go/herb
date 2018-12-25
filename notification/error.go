package notification

import "strings"

type NotificationError struct {
	Instance *NotificationInstance
	Err      error
}

func (e *NotificationError) Error() string {
	output := "notification error \n"
	output += "notificaion error message :" + e.Err.Error() + "\n"
	output += "notificaion intstance ID:" + e.Instance.InstanceID + "\n"
	output += "notificaion intstance Sender:" + e.Instance.Sender + "\n"
	output += "notificaion Recipient :" + e.Instance.Recipient + "\n"
	output += "notificaion Logs :" + strings.Join(e.Instance.Logs, "\n") + "\n"
	return output
}
