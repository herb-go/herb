package credential

type Credential interface {
	Type() Type
	Data() ([]byte, error)
}

type PlainCredential struct {
	credentialType Type
	credentialData []byte
}

func (c *PlainCredential) Type() Type {
	return c.credentialType
}
func (c *PlainCredential) WithType(t Type) *PlainCredential {
	c.credentialType = t
	return c
}
func (c *PlainCredential) Data() ([]byte, error) {
	return c.credentialData, nil
}

func (c *PlainCredential) WithData(data []byte) *PlainCredential {
	c.credentialData = data
	return c
}
func New() *PlainCredential {
	return &PlainCredential{}
}

type Type string

var TypeUID = Type("uid")
var TypePassword = Type("password")
var TypeAppID = Type("appid")
var TypeToken = Type("token")
var TypeTimestamp = Type("timestamp")
var TypeSign = Type("sign")
var TypeSession = Type("session")

type Collection map[Type][]byte

func (c *Collection) Set(t Type, v []byte) {
	(*c)[t] = v
}
func (c *Collection) Get(t Type) []byte {
	return (*c)[t]
}

func NewCollection() *Collection {
	return &Collection{}
}
