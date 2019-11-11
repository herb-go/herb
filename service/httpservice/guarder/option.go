package guarder

import "github.com/herb-go/herb/server"

type GuarderOption interface {
	ApplyToGuarder(g *Guarder) error
}

type VisitorOption interface {
	ApplyToVisitor(v *Visitor) error
}

type DriverConfig interface {
	MapperDriverName() string
	DriverName() string
	DriverConfig() server.Config
}

func ApplyToGuarder(g *Guarder, c DriverConfig) error {
	config := c.DriverConfig()
	if g.Mapper == nil {
		d := c.MapperDriverName()
		driver, err := NewMapperDriver(d, config, "")
		if err != nil {
			return err
		}
		g.Mapper = driver
	}
	if g.Identifier == nil {
		d := c.DriverName()
		driver, err := NewIdentifierDriver(d, config, "")
		if err != nil {
			return err
		}
		g.Identifier = driver

	}
	return nil

}

func ApplyToVisitor(v *Visitor, c DriverConfig) error {
	config := c.DriverConfig()
	if v.Mapper == nil {
		d := c.MapperDriverName()
		driver, err := NewMapperDriver(d, config, "")
		if err != nil {
			return err
		}
		v.Mapper = driver
	}
	if v.Credential == nil {
		d := c.DriverName()
		driver, err := NewCredentialDriver(d, config, "")
		if err != nil {
			return err
		}
		v.Credential = driver

	}
	return nil

}

type DriverField struct {
	Driver       string
	staticDriver string
}

func (f *DriverField) SetStaticDriver(d string) {
	f.staticDriver = d
}
func (f *DriverField) DriverName() string {
	if f.staticDriver == "" {
		return f.Driver
	}
	return f.staticDriver
}

type MapperDriverField struct {
	MapperDriver       string
	staticMapperDriver string
}

func (f *MapperDriverField) SetStaticMapperDriver(d string) {
	f.staticMapperDriver = d
}
func (f *MapperDriverField) MapperDriverName() string {
	if f.staticMapperDriver == "" {
		return f.MapperDriver
	}
	return f.staticMapperDriver
}

func NewDriverConfigMap() *DirverConfigMap {
	return &DirverConfigMap{}
}

type DirverConfigMap struct {
	DriverField
	MapperDriverField
	Config server.ConfigMap
}

func (c *DirverConfigMap) DriverConfig() server.Config {
	return &c.Config
}

func (c *DirverConfigMap) ApplyToGuarder(g *Guarder) error {
	return ApplyToGuarder(g, c)
}

func (c *DirverConfigMap) ApplyToVisitor(v *Visitor) error {
	return ApplyToVisitor(v, c)
}
