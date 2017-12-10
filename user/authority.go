package user

import "net/http"

type AuthorityChecker interface {
	Check(user string, authority string) (bool, error)
}

func NewAuthorityService(identifier Identifier, checker AuthorityChecker) *AuthorityService {
	return &AuthorityService{
		identifier: identifier,
		checker:    checker,
	}
}

type AuthorityService struct {
	identifier Identifier
	checker    AuthorityChecker
}

func (m *AuthorityService) CheckAuthorityMiddleware(authority string, unauthorizedAction http.HandlerFunc) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		id, err := m.identifier.IdentifyRequest(r)
		if err != nil {
			panic(err)
		}
		result, err := m.checker.Check(id, authority)
		if err != nil {
			panic(err)
		}
		if result != true {
			if unauthorizedAction == nil {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			} else {
				unauthorizedAction(w, r)
			}
			return
		}
		next(w, r)
	}
}

func (m *AuthorityService) CheckAuthoritiesMiddleware(authorities []string, unauthorizedAction http.HandlerFunc) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		id, err := m.identifier.IdentifyRequest(r)
		if err != nil {
			panic(err)
		}
		result, err := m.CheckMutliAuthorities(id, authorities)
		if err != nil {
			panic(err)
		}
		if result != true {
			if unauthorizedAction == nil {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			} else {
				unauthorizedAction(w, r)
			}
			return
		}
		next(w, r)
	}
}
func (m *AuthorityService) CheckAuthoritiesOrForbiddenMiddleware(authorities ...string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		id, err := m.identifier.IdentifyRequest(r)
		if err != nil {
			panic(err)
		}
		result, err := m.CheckMutliAuthorities(id, authorities)
		if err != nil {
			panic(err)
		}
		if result != true {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

func (m *AuthorityService) CheckMutliAuthorities(id string, authorities []string) (bool, error) {
	if len(authorities) == 0 {
		return false, nil
	}
	for _, v := range authorities {
		result, err := m.checker.Check(id, v)
		if err != nil {
			return false, err
		}
		if result == true {
			return true, nil
		}
	}
	return false, nil
}

func (m *AuthorityService) Check(user string, authority string) (bool, error) {
	return m.checker.Check(user, authority)
}
