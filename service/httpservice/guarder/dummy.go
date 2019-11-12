package guarder

import (
	"github.com/herb-go/herb/service"
)

func NewDummy() *Dummy {
	return &Dummy{}
}

type Dummy struct {
	ID string
}

func (d *Dummy) IdentifyParams(p *Params) (string, error) {
	return d.ID, nil
}

func blockAllIdentifierFactory(conf service.Config, prefix string) (Identifier, error) {
	d := NewDummy()
	return d, nil
}
func dummyIdentifierFactory(conf service.Config, prefix string) (Identifier, error) {
	var err error
	d := NewDummy()
	err = conf.Get("ID", &d.ID)
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
