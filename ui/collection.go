package ui

//Collection struct which hold avaliable message map and translated messages.
type Collection struct {
	messages    *Messages
	messagesmap map[string]string
}

//NewCollection create new map with given avaliable map and translated messages.
func NewCollection(messages *Messages, messagesmap map[string]string) *Collection {
	return &Collection{
		messages:    messages,
		messagesmap: messagesmap,
	}
}

//Get get translated field label .
func (m *Collection) Get(field string) string {
	label := m.messagesmap[field]
	if label == "" {
		return field
	}
	if m.messages == nil {
		return label
	}
	return m.messages.Get(label)
}
