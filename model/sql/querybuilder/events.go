package querybuilder

type CommonQueryEvents struct {
}

func (e CommonQueryEvents) BeforeInsert() error {
	return nil
}

func (e CommonQueryEvents) BeforeUpdate() error {
	return nil
}

func (e CommonQueryEvents) AfterFind() error {
	return nil
}

func (e CommonQueryEvents) AfterDelete() error {
	return nil
}
