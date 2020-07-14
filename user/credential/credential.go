package credential

type Credential interface {
	Type() Type
	Load() (Data, error)
}

type Type string

var TypeUID = Type("uid")
var TypePassword = Type("password")
var TypeAppID = Type("appid")
var TypeToken = Type("token")
var TypeTimestamp = Type("timestamp")
var TypeSign = Type("sign")
var TypeSession = Type("session")

type Data []byte

type Collection map[Type]Data

func (c *Collection) Set(t Type, v Data) {
	(*c)[t] = v
}
func (c *Collection) Get(t Type) Data {
	return (*c)[t]
}

func NewCollection() *Collection {
	return &Collection{}
}
