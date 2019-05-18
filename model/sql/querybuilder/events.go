package querybuilder

//CommonQueryEvents common query events struct
type CommonQueryEvents struct {
}

//BeforeInsert event raised before insert
func (e CommonQueryEvents) BeforeInsert() error {
	return nil
}

//BeforeInsert event raised after insert
func (e CommonQueryEvents) AfterInsert() error {
	return nil
}

//BeforeInsert event raised before update
func (e CommonQueryEvents) BeforeUpdate() error {
	return nil
}

//BeforeInsert event raised after update
func (e CommonQueryEvents) AfterUpdate() error {
	return nil
}

//BeforeInsert event raised before find
func (e CommonQueryEvents) AfterFind() error {
	return nil
}

//BeforeInsert event raised after find
func (e CommonQueryEvents) AfterDelete() error {
	return nil
}
