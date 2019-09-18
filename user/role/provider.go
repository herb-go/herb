package role

//Provider roles provider interface
type Provider interface {
	//Roles get roles by user id.
	//Return user roles and any error if raised.
	Roles(uid string) (*Roles, error)
}
