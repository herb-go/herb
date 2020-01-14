package requestmatching

import "net/http"

type Pattern interface {
	MatchRequest(r *http.Request) (bool, error)
}

func MatchAll(r *http.Request, f ...Pattern) (bool, error) {
	var result bool
	var err error
	for k := range f {
		result, err = f[k].MatchRequest(r)
		if err != nil || result == false {
			return false, err
		}
	}
	return true, nil
}

func MatchAny(r *http.Request, f ...Pattern) (bool, error) {
	var result bool
	var err error
	for k := range f {
		result, err = f[k].MatchRequest(r)
		if err != nil {
			return false, err
		}
		if result {
			return true, nil
		}
	}
	return false, nil
}

type PlainPattern struct {
	IPNets  IPNets
	Methods Methods
	Exts    Exts
	Paths   PathsWithHost
	Prefixs PrefixsWithHost
}

func (p *PlainPattern) MatchRequest(r *http.Request) (bool, error) {
	return MatchAll(r,
		p.IPNets,
		p.Methods,
		p.Exts,
		p.Paths,
		p.Prefixs,
	)
}
