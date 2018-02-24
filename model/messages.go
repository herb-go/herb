package model

//Messages translate messages map.
type Messages map[string]string

//GetMessage get translated string for key.
//Return key if translateed string not exist.
func (m Messages) GetMessage(key string) string {
	result, ok := m[key]
	if ok == false {
		result = key
	}
	return result
}

//HasMessage check if translated string exists for given key.
//If string exists,return tranlasted string and true.
//If string does not exist,return key and false.
func (m Messages) HasMessage(key string) (string, bool) {
	value, ok := m[key]
	if ok == false {
		value = key
	}
	return value, ok
}

//MessageChain use model messages interface list as model messages interface.
type MessageChain []MessagesCollection

//GetMessage get translated string for key.
//Return key if translateed string not exist.
//Check all model messages in order.
func (m *MessageChain) GetMessage(key string) string {
	if m != nil {
		for _, v := range *m {
			value, ok := v.HasMessage(key)
			if ok {
				return value
			}
		}
	}
	return key
}

//HasMessage check if translated string exists for given key.
//If string exists,return tranlasted string and true.
//If string does not exist,return key and false.
//Check all model messages in order.
func (m *MessageChain) HasMessage(key string) (string, bool) {
	if m != nil {
		for _, v := range *m {
			value, ok := v.HasMessage(key)
			if ok {
				return value, ok
			}
		}
	}
	return key, false
}

//Use append new model messages to MessageChain.
func (m *MessageChain) Use(Messages ...MessagesCollection) *MessageChain {
	*m = append(*m, Messages...)
	return m
}

//NewMessageChain create new message chain with given model messages.
func NewMessageChain(Messages ...MessagesCollection) *MessageChain {
	m := make([]MessagesCollection, len(Messages))
	copy(m, Messages)
	chain := MessageChain(m)
	return &chain
}

//DefaultMessages default messages
var DefaultMessages MessageChain

//MessagesCollection model messages collection interface.
type MessagesCollection interface {
	//GetMessage get translated string for key.
	//Return key if translateed string not exist.
	GetMessage(key string) string
	//HasMessage check if translated string exists for given key.
	//If string exists,return tranlasted string and true.
	//If string does not exist,return key and false.
	HasMessage(key string) (string, bool)
}
