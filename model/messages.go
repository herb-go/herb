package model

type Messages map[string]string

func (m Messages) GetMessage(key string) string {
	result, ok := m[key]
	if ok == false {
		result = key
	}
	return result
}
func (m Messages) HasMessage(key string) (string, bool) {
	value, ok := m[key]
	if ok == false {
		value = key
	}
	return value, ok
}

type MessageChain []ModelMessages

func (m *MessageChain) GetMessage(key string) string {
	for _, v := range *m {
		value, ok := v.HasMessage(key)
		if ok {
			return value
		}
	}
	return key
}
func (m *MessageChain) HasMessage(key string) (string, bool) {
	for _, v := range *m {
		value, ok := v.HasMessage(key)
		if ok {
			return value, ok
		}
	}
	return key, false
}

func (m *MessageChain) Use(Messages ...ModelMessages) *MessageChain {
	messageLength := len(Messages)
	backup := make([]ModelMessages, len(*m))
	copy(backup, *m)
	*m = make([]ModelMessages, len(*m)+messageLength)
	copy((*m)[0:messageLength], Messages)
	for i := 0; i < messageLength; i++ {
		(*m)[i] = Messages[messageLength-1-i]
	}
	return m
}

func NewMessageChain(Messages ...ModelMessages) *MessageChain {
	m := make([]ModelMessages, len(Messages))
	copy(m, Messages)
	chain := MessageChain(m)
	return &chain
}

var DefaultMessages = NewMessageChain()

type ModelMessages interface {
	GetMessage(key string) string
	HasMessage(key string) (string, bool)
}
