package role

import (
	"net/http"

	"github.com/herb-go/herb/user"
)

//RuleProvider rule provider interface
type RuleProvider interface {
	//Rule get rule from http request.
	//Return rule and any error if raised.
	Rule(*http.Request) (Rule, error)
}

//Provider roles provider interface
type Provider interface {
	//Roles get roles by user id.
	//Return user roles and any error if raised.
	Roles(uid string) (*Roles, error)
}

//Authorizer role service authorizer
type Authorizer struct {
	Service      *Service
	RuleProvider RuleProvider
}

//Authorize authorized with requestWW
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

//NewService create new role authorize service with role provider and user Identifier.
func NewService(RoleProvider Provider, Identifier user.Identifier) *Service {
	return &Service{
		RoleProvider: RoleProvider,
		Identifier:   Identifier,
	}
}

//Service role authorize service
type Service struct {
	RoleProvider Provider
	Identifier   user.Identifier
}

//Authorizer create authorizer with given rule provider
func (s *Service) Authorizer(rs RuleProvider) *Authorizer {
	return &Authorizer{
		Service:      s,
		RuleProvider: rs,
	}
}

//AuthorizeMiddleware middleware which authorize http request.
//If authorize fail,unauthorizedAction will be executed.
func (s *Service) AuthorizeMiddleware(rs RuleProvider, unauthorizedAction http.HandlerFunc) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return user.AuthorizeMiddleware(s.Authorizer(rs), unauthorizedAction)
}

//RolesAuthorizeOrForbiddenMiddleware middleware which authorize http request.
//If authorize fail,http error forbidden will be executed.
func (s *Service) RolesAuthorizeOrForbiddenMiddleware(ruleNames ...string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var rs = NewRoles(ruleNames...)
	return s.AuthorizeMiddleware(rs, nil)
}
