package application

import "github.com/herb-go/herbsecurity/authority"

type Creator interface {
	CreateApplication(authority.Principal, authority.Agent) Verified
}

type Regenerator interface {
	RegenerateApplication(authority.Principal, authority.Authority)
}

type Revoker interface {
	RevokeApplication(authority.Principal, authority.Authority)
}
type Loader interface {
	LoadApplication(authority.Authority) (*Verified, error)
}
