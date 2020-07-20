package credential

import "github.com/herb-go/herbsecurity/authority"

type Authenticator interface {
	Authenticate(Credentials) (*authority.Authentication, error)
}

type Dependences map[Name]bool

func (d Dependences) DependencesData() (Dependences, error) {
	return d, nil
}

type DependencesSource interface {
	DependencesData() (Dependences, error)
}

type DependencesAuthenticator interface {
	Authenticator
	DependencesSource
}
