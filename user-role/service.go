package role

import (
	"net/http"

	"github.com/herb-go/herb/user"
)

type RuleProvider interface {
	Rule(*http.Request) (Rule, error)
}
type RoleProvider interface {
	Roles(uid string) (*Roles, error)
}

type Authorizer struct {
	Service      *Service
	RuleProvider RuleProvider
}

func (a *Authorizer) Authorize(r *http.Request) (bool, error) {
	uid, err := a.Service.Identifier.IdentifyRequest(r)
	if err != nil {
		return false, err
	}
	if uid == "" {
		return false, nil
	}
	roles, err := a.Service.RoleProvider.Roles(uid)
	if err != nil {
		return false, err
	}
	if roles == nil {
		return false, err
	}
	rm, err := a.RuleProvider.Rule(r)
	if err != nil {
		return false, err
	}
	return rm.Execute(*roles...)
}

type Service struct {
	RoleProvider RoleProvider
	Identifier   user.Identifier
}

func (s *Service) Authorizer(rs RuleProvider) *Authorizer {
	return &Authorizer{
		Service:      s,
		RuleProvider: rs,
	}
}
func (s *Service) AuthorizeMiddleware(rs RuleProvider, unauthorizedAction http.HandlerFunc) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return user.AuthorizeMiddleware(s.Authorizer(rs), unauthorizedAction)
}

func (s *Service) RolesAuthorizeOrForbiddenMiddleware(ruleNames ...string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var rs = NewRoles(ruleNames...)
	return s.AuthorizeMiddleware(rs, nil)
}
