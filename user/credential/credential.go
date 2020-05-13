package credential

type Loader interface {
	CredentialType() Type
	LoadCredential() (interface{}, error)
}

type FixedLoader struct {
	value interface{}
}

func (l *FixedLoader) LoadCredential() (interface{}, error) {
	return l.value, nil
}

func NewFixedLoader(v interface{}) *FixedLoader {
	return &FixedLoader{
		value: v,
	}
}

type Type string

type Credential struct {
	Type  Type
	Value interface{}
}

func New() *Credential {
	return &Credential{}
}

type Data map[Type][]interface{}

func (d *Data) Append(t Type, v interface{}) {
	(*d)[t] = append((*d)[t], v)
}
func (d *Data) Get(t Type) interface{} {
	if len((*d)[t]) == 0 {
		return nil
	}
	return (*d)[t][0]
}
func (d *Data) GetAllByType(t Type) []interface{} {
	return (*d)[t]
}
func NewData() *Data {
	return &Data{}
}
