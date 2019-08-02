package ui

//Messages translate messages map.
type Messages map[string]string

//NewMessages create new messages
func NewMessages() *Messages {
	return &Messages{}
}

//Get get translated string for key.
//Return key if translateed string not exist.
func (m *Messages) Get(key string) string {
	result, ok := (*m)[key]
	if ok == false {
		result = key
	}
	return result
}

//Set set translated string to key.
//Return messages self.
func (m *Messages) Set(key string, message string) *Messages {
	(*m)[key] = message
	return m
}

//Collection create collection with given messagesmap
func (m *Messages) Collection(messagesmap map[string]string) *Collection {
	return NewCollection(m, messagesmap)
}

//Load check if translated string exists for given key.
//If string exists,return tranlasted string and true.
//If string does not exist,return key and false.
func (m *Messages) Load(key string) (string, bool) {
	value, ok := (*m)[key]
	if ok == false {
		value = key
	}
	return value, ok
}
