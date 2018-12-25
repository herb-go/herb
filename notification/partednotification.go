package notification

type PartedNotification struct {
	CommonNotification
	Parts map[string][]byte
}

func NewPartedNotification() *PartedNotification {
	return &PartedNotification{
		Parts: map[string][]byte{},
	}
}

func NewPartedNotificationWithID() (*PartedNotification, error) {
	id, err := DefaultIDGenerator()
	if err != nil {
		return nil, err
	}
	n := NewPartedNotification()
	n.SetNotificationID(id)
	return n, nil
}

func ValidatePartedNotificationInstanceWithFields(ni *NotificationInstance, fields [][]string) (*PartedNotification, error) {
	pn, ok := ni.Notification.(*PartedNotification)
	if ok == false {
		ni.SetStatusUnsupported()
		return nil, nil
	}
	if fields != nil {
		for mustIndex := range fields {
			result := false
			for oneofIndex := range fields[mustIndex] {
				field := fields[mustIndex][oneofIndex]
				part := pn.Parts[field]
				if len(part) > 0 {
					result = true
					break
				}
			}
			if result == false {
				ni.SetStatusUnsupported()
				return nil, nil
			}
		}
	}
	return pn, nil
}
