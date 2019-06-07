package model

//Messages translate messages map.
type Messages map[string]string

//NewMessages create new messages
func NewMessages() *Messages {
	return &Messages{}
}

//GetMessage get translated string for key.
//Return key if translateed string not exist.
func (m *Messages) GetMessage(key string) string {
	result, ok := (*m)[key]
	if ok == false {
		result = key
	}
	return result
}

//SetMessage set translated string to key.
//Return messages self.
func (m *Messages) SetMessage(key string, message string) *Messages {
	(*m)[key] = message
	return m
}

//LoadMessage check if translated string exists for given key.
//If string exists,return tranlasted string and true.
//If string does not exist,return key and false.
func (m *Messages) LoadMessage(key string) (string, bool) {
	value, ok := (*m)[key]
	if ok == false {
		value = key
	}
	return value, ok
}

//MessagesChain use model messages interface list as model messages interface.
type MessagesChain []MessagesCollection

//GetMessage get translated string for key.
//Return key if translateed string not exist.
//Check all model messages in order.
func (m *MessagesChain) GetMessage(key string) string {
	if m != nil {
		for _, v := range *m {
			value, ok := v.LoadMessage(key)
			if ok {
				return value
			}
		}
	}
	return key
}

//LoadMessage check if translated string exists for given key.
//If string exists,return tranlasted string and true.
//If string does not exist,return key and false.
//Check all model messages in order.
func (m *MessagesChain) LoadMessage(key string) (string, bool) {
	if m != nil {
		for _, v := range *m {
			value, ok := v.LoadMessage(key)
			if ok {
				return value, ok
			}
		}
	}
	return key, false
}

//Use append new model messages to MessagesChain.
func (m *MessagesChain) Use(Messages ...MessagesCollection) *MessagesChain {
	*m = append(*m, Messages...)
	return m
}

//NewMessagesChain create new message chain with given model messages.
func NewMessagesChain(Messages ...MessagesCollection) *MessagesChain {
	m := make([]MessagesCollection, len(Messages))
	copy(m, Messages)
	chain := MessagesChain(m)
	return &chain
}

//DefaultMessagesChain default messages
var DefaultMessagesChain = NewMessagesChain()

//MessagesCollection model messages collection interface.
type MessagesCollection interface {
	//GetMessage get translated string for key.
	//Return key if translateed string not exist.
	GetMessage(key string) string
	//LoadMessage check if translated string exists for given key.
	//If string exists,return tranlasted string and true.
	//If string does not exist,return key and false.
	LoadMessage(key string) (string, bool)
}

//Use append new model messages to default MessagesChain.
func Use(Messages ...MessagesCollection) *MessagesChain {
	return DefaultMessagesChain.Use(Messages...)
}

//GetMessage get translated string for key from default  MessagesChain.
//Return key if translateed string not exist.
func GetMessage(key string) string {
	return DefaultMessagesChain.GetMessage(key)
}
