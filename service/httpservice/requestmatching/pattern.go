package requestmatching

import (
	"net/http"
)

//MustMatch match request with given pattern
//Return false if pattern is nil.
//Panic if any error raised.
func MustMatch(r *http.Request, p Pattern) bool {
	if p == nil {
		return false
	}
	result, err := p.MatchRequest(r)
	if err != nil {
		panic(err)
	}
	return result
}

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

//MustMatchAll match given request and given partterns.
//Match will fail if any pattern fail.
//Return result and painc if any error raised.
func MustMatchAll(r *http.Request, p ...Pattern) bool {
	result, err := MatchAll(r, p...)
	if err != nil {
		panic(err)
	}
	return result
}

//MatchAny match given request and given partterns.
//Match will success if any pattern success or no patterns given.
//Return result and any error if raised.
func MatchAny(r *http.Request, p ...Pattern) (bool, error) {
	var result bool
	var err error
	if len(p) == 0 {
		return true, nil
	}
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

//MustMatchAny match given request and given partterns.
//Match will success if any pattern success.
//Return result and painc if any error raised.
func MustMatchAny(r *http.Request, p ...Pattern) bool {
	result, err := MatchAny(r, p...)
	if err != nil {
		panic(err)
	}
	return result
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

//Filters pattern filters type
type Filters []Pattern

//MatchRequest match request.
//Return result and any error if raised.
func (f *Filters) MatchRequest(r *http.Request) (bool, error) {
	return MatchAny(r, (*f)...)
}

//FiltersConfig filters config struct
type FiltersConfig []*PatternConfig

//CreatePattern create plain pattern.
//Return pattern created and any error if raised.
//Match will success if filters is empty
func (c *FiltersConfig) CreatePattern() (Pattern, error) {
	f := Filters{}
	for k := range *c {
		pattern, err := (*c)[k].CreatePattern()
		if err != nil {
			return nil, err
		}
		f = append(f, pattern)
	}
	return &f, nil
}

//Whitelist whitelist pattern type
type Whitelist []Pattern

//MatchRequest match request.
//Return result and any error if raised.
//Match will fail if whitelist is empty
func (w *Whitelist) MatchRequest(r *http.Request) (bool, error) {
	if len(*w) == 0 {
		return false, nil
	}
	return MatchAny(r, (*w)...)
}

//WhitelistConfig whitelist config type
type WhitelistConfig []*PatternConfig

//CreatePattern create plain pattern.
//Return pattern created and any error if raised.
func (c *WhitelistConfig) CreatePattern() (Pattern, error) {
	w := Whitelist{}
	for k := range *c {
		pattern, err := (*c)[k].CreatePattern()
		if err != nil {
			return nil, err
		}
		w = append(w, pattern)
	}
	return &w, nil
}

//PatternAll pattern all type
type PatternAll []Pattern

//MatchRequest match request.
//Return result and any error if raised.
func (p *PatternAll) MatchRequest(r *http.Request) (bool, error) {
	return MatchAll(r, (*p)...)
}

//PatternAllConfig pattern all type
type PatternAllConfig []*PatternConfig

//CreatePattern create plain pattern.
//Return pattern created and any error if raised.
func (c *PatternAllConfig) CreatePattern() (Pattern, error) {
	p := PatternAll{}
	for k := range *c {
		pattern, err := (*c)[k].CreatePattern()
		if err != nil {
			return nil, err
		}
		p = append(p, pattern)
	}
	return &p, nil
}

//Factory Factory interface
type Factory interface {
	//CreatePattern create plain pattern.
	//Retyurn pattern created and any error if raised.
	CreatePattern() (Pattern, error)
}

//MustCreatePattern create pattern with given factory.
//Panic if any error raised.
func MustCreatePattern(f Factory) Pattern {
	p, err := f.CreatePattern()
	if err != nil {
		panic(err)
	}
	return p
}
