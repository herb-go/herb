package requestmatching

import (
	"context"
	"net"
	"net/http"
)

//GetRequestIPAddress get ip address from given request.
func GetRequestIPAddress(r *http.Request) string {
	v := r.Context().Value(ContextKeyIPAddress)
	if v != nil {
		return v.(string)
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	ctx := context.WithValue(r.Context(), ContextKeyIPAddress, ip)
	req := r.WithContext(ctx)
	*r = *req
	return ip

}

//GetRequestIP get ip from given request.
func GetRequestIP(r *http.Request) net.IP {
	v := r.Context().Value(ContextKeyIP)
	if v != nil {
		return v.(net.IP)
	}
	ip := net.ParseIP(GetRequestIPAddress(r))
	ctx := context.WithValue(r.Context(), ContextKeyIP, ip)
	req := r.WithContext(ctx)
	*r = *req
	return ip
}

//IPNets ip nets pattern.
type IPNets []*net.IPNet

//MatchRequest match request.
//Return result and any error if raised
func (i *IPNets) MatchRequest(r *http.Request) (bool, error) {
	if len(*i) == 0 {
		return true, nil
	}
	ip := GetRequestIP(r)
	if ip == nil {
		return false, nil
	}
	for k := range *i {
		if (*i)[k].Contains(ip) {
			return true, nil
		}
	}
	return false, nil
}

//Add add CIDR format ip net to pattern.
//Return any error if raised.
func (i *IPNets) Add(pattern string) error {
	_, ipnet, err := net.ParseCIDR(pattern)
	if err != nil {
		return err
	}
	*i = append(*i, ipnet)
	return nil
}

//NewIPNets new ip nets pattern.
func NewIPNets() *IPNets {
	return &IPNets{}
}
