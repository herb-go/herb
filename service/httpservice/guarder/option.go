package guarder

type GuarderOption interface {
	ApplyToGuarder(g *Guarder) error
}

type VisitorOption interface {
	ApplyToVisitor(v *Visitor) error
}

func NewDriverConfig() *DriverConfig {
	return &DriverConfig{}
}

type DriverConfig struct {
	DriverField
	MapperDriverField
	Config func(interface{}) error `config:", lazyload"`
}

func (c *DriverConfig) ApplyToGuarder(g *Guarder) error {
	return ApplyToGuarder(g, c)
}

func (c *DriverConfig) ApplyToVisitor(v *Visitor) error {
	return ApplyToVisitor(v, c)
}

func ApplyToGuarder(g *Guarder, c *DriverConfig) error {
	if g.Mapper == nil {
		d := c.MapperDriverName()
		driver, err := NewMapperDriver(d, c.Config)
		if err != nil {
			return err
		}
		g.Mapper = driver
	}
	if g.Identifier == nil {
		d := c.DriverName()
		driver, err := NewIdentifierDriver(d, c.Config)
		if err != nil {
			return err
		}
		g.Identifier = driver

	}
	return nil

}

func ApplyToVisitor(v *Visitor, c *DriverConfig) error {
	if v.Mapper == nil {
		d := c.MapperDriverName()
		driver, err := NewMapperDriver(d, c.Config)
		if err != nil {
			return err
		}
		v.Mapper = driver
	}
	if v.Credential == nil {
		d := c.DriverName()
		driver, err := NewCredentialDriver(d, c.Config)
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
