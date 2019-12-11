package guarder

func NewDummy() *Dummy {
	return &Dummy{}
}

type Dummy struct {
	ID string
}

func (d *Dummy) IdentifyParams(p *Params) (string, error) {
	return d.ID, nil
}

func blockAllIdentifierFactory(loader func(interface{}) error) (Identifier, error) {
	d := NewDummy()
	return d, nil
}
func dummyIdentifierFactory(loader func(interface{}) error) (Identifier, error) {
	var err error
	d := NewDummy()
	err = loader(d)
	if err != nil {
		return nil, err
	}

	return d, nil
}
func registerDummyFactory() {
	RegisterIdentifier("", blockAllIdentifierFactory)
	RegisterIdentifier("dummy", dummyIdentifierFactory)
}

func init() {
	registerDummyFactory()
}
