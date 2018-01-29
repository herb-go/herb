package role

import (
	"net/http"

	"github.com/herb-go/herb/user"
)

type RuleService interface {
	Rule(*http.Request) (Rule, error)
}
type RoleService interface {
	Roles(uid string) (*Roles, error)
}

type Authorizer struct {
	Service     *Service
	RuleService RuleService
}

func (a *Authorizer) Authorize(r *http.Request) (bool, error) {
	uid, err := a.Service.Identifier.IdentifyRequest(r)
	if err != nil {
		return false, err
	}
	if uid == "" {
		return false, nil
	}
	roles, err := a.Service.RoleService.Roles(uid)
	if err != nil {
		return false, err
	}
	if roles == nil {
		return false, err
	}
	rm, err := a.RuleService.Rule(r)
	if err != nil {
		return false, err
	}
	return rm.Execute(*roles...)
}

type Service struct {
	RoleService RoleService
	Identifier  user.Identifier
}

func (s *Service) Authorizer(rs RuleService) *Authorizer {
	return &Authorizer{
		Service:     s,
		RuleService: rs,
	}
}
func (s *Service) AuthorizeMiddleware(rs RuleService, unauthorizedAction http.HandlerFunc) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return user.AuthorizeMiddleware(s.Authorizer(rs), unauthorizedAction)
}

func (s *Service) RolesAuthorizeOrForbiddenMiddleware(ruleNames ...string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var rs = NewRoles(ruleNames...)
	return s.AuthorizeMiddleware(rs, nil)
}
