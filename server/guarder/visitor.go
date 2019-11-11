package guarder

import (
	"net/http"
)

type Visitor struct {
	Credential Credential
	Mapper     Mapper
}

func (c *Visitor) CredentialRequest(r *http.Request) error {
	p, err := c.Credential.CredentialParams()
	if err != nil {
		return err
	}
	return c.Mapper.WriteParamsToRequest(r, p)
}

func (c *Visitor) Init(o VisitorOption) error {
	return o.ApplyToVisitor(c)
}
func NewVisitor() *Visitor {
	return &Visitor{}
}
