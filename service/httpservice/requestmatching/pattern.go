package requestmatching

import (
	"net/http"
)

//Pattern request matching pattern interafe
type Pattern interface {
	//MatchRequest match request.
	//Return result and any error if raised.
	MatchRequest(r *http.Request) (bool, error)
}

//MatchAll match given request and given partterns.
//Match will fail if any pattern fail.
//Return result and any error if raised.
func MatchAll(r *http.Request, p ...Pattern) (bool, error) {
	var result bool
	var err error
	for k := range p {
		result, err = p[k].MatchRequest(r)
		if err != nil || result == false {
			return false, err
		}
	}
	return true, nil
}

//MatchAny match given request and given partterns.
//Match will success if any pattern success.
//Return result and any error if raised.
func MatchAny(r *http.Request, p ...Pattern) (bool, error) {
	var result bool
	var err error
	for k := range p {
		result, err = p[k].MatchRequest(r)
		if err != nil {
			return false, err
		}
		if result {
			return true, nil
		}
	}
	return false, nil
}

//PlainPattern plain pattern struct
type PlainPattern struct {
	IPNets   *IPNets
	Methods  *Methods
	Exts     *Exts
	Paths    *Paths
	Prefixs  *Prefixs
	Disabled bool
	Not      bool
	And      bool
	Patterns []Pattern
}

//NewPlainPattern create new pattern struct.
func NewPlainPattern() *PlainPattern {
	return &PlainPattern{
		IPNets:  NewIPNets(),
		Methods: NewMethods(),
		Exts:    NewExts(),
		Paths:   NewPaths(),
		Prefixs: NewPrefixs(),
	}
}

//MatchRequest match request.
//Return result and any error if raised.
func (p *PlainPattern) MatchRequest(r *http.Request) (bool, error) {
	if p.Disabled {
		return false, nil
	}
	result, err := p.matchRequest(r)
	if err != nil {
		return false, err
	}
	return result != p.Not, nil
}

func (p *PlainPattern) matchRequest(r *http.Request) (bool, error) {
	result, err := MatchAll(r,
		p.IPNets,
		p.Methods,
		p.Exts,
		p.Paths,
		p.Prefixs,
	)
	if err != nil {
		return false, err
	}
	if p.Patterns == nil {
		return result, nil
	}
	if p.And {
		return MatchAll(r, p.Patterns...)
	}
	if result == true {
		return true, nil
	}
	return MatchAny(r, p.Patterns...)
}

//PatternConfig plainpattern config struct
type PatternConfig struct {
	IPList     []string
	URLList    []string
	PrefixList []string
	ExtList    []string
	MethodList []string
	Disabled   bool
	Not        bool
	And        bool
	Patterns   []*PatternConfig
}

//CreatePattern create plain pattern.
//Retyurn pattern created and any error if raised.
func (c *PatternConfig) CreatePattern() (Pattern, error) {
	p := NewPlainPattern()
	for k := range c.IPList {
		err := p.IPNets.Add(c.IPList[k])
		if err != nil {
			return nil, err
		}
	}
	for k := range c.URLList {
		p.Paths.Add(c.URLList[k])
	}
	for k := range c.PrefixList {
		p.Prefixs.Add(c.PrefixList[k])
	}
	for k := range c.ExtList {
		p.Exts.Add(c.ExtList[k])
	}
	for k := range c.MethodList {
		p.Methods.Add(c.MethodList[k])
	}
	p.Disabled = c.Disabled
	p.Not = c.Not
	p.And = c.And
	for k := range c.Patterns {
		pattern, err := c.Patterns[k].CreatePattern()
		if err != nil {
			return nil, err
		}
		p.Patterns = append(p.Patterns, pattern)
	}
	return p, nil
}
