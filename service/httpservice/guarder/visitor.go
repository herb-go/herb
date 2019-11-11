package guarder

import (
	"net/http"
)

type Visitor struct {
	Credential Credential
	Mapper     Mapper
}

func (v *Visitor) InitRequest(r *http.Request) error {
	return v.CredentialRequest(r)
}

func (v *Visitor) CredentialRequest(r *http.Request) error {
	p, err := v.Credential.CredentialParams()
	if err != nil {
		return err
	}
	return v.Mapper.WriteParamsToRequest(r, p)
}

func (v *Visitor) Init(o VisitorOption) error {
	return o.ApplyToVisitor(v)
}
func NewVisitor() *Visitor {
	return &Visitor{}
}
