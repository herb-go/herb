package roleservice

import (
	"net/http"

	"github.com/herb-go/herb/user/httpuser"
	"github.com/herb-go/herb/user/role"
)

//RuleProvider rule provider interface
type RuleProvider interface {
	//Rule get rule from http request.
	//Return rule and any error if raised.
	Rule(*http.Request) (role.Rule, error)
}

//Authorizer role service authorizer
type Authorizer struct {
	Service      *Service
	RuleProvider RuleProvider
}

//RolesFromRequest get roles from request
func (a *Authorizer) RolesFromRequest(r *http.Request) (*role.Roles, error) {
	uid, err := a.Service.Identifier.IdentifyRequest(r)
	if err != nil {
		return nil, err
	}
	if uid == "" {
		return nil, nil
	}
	return a.Service.RoleProvider.Roles(uid)
}

//Authorize authorized with requestWW
func (a *Authorizer) Authorize(r *http.Request) (bool, error) {
	roles, err := a.RolesFromRequest(r)
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
func NewService(RoleProvider role.Provider, Identifier httpuser.Identifier) *Service {
	return &Service{
		RoleProvider: RoleProvider,
		Identifier:   Identifier,
	}
}

//Service role authorize service
type Service struct {
	RoleProvider role.Provider
	Identifier   httpuser.Identifier
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
	return httpuser.AuthorizeMiddleware(s.Authorizer(rs), unauthorizedAction)
}

//RolesAuthorizeOrForbiddenMiddleware middleware which authorize http request.
//If authorize fail,http error forbidden will be executed.
func (s *Service) RolesAuthorizeOrForbiddenMiddleware(ruleNames ...string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var rs = role.NewRoles(ruleNames...)
	return s.AuthorizeMiddleware(rs, nil)
}
