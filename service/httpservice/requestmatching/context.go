package requestmatching

//ContextKey context key type
type ContextKey string

//ContextKeyIPAddress context key for request ip address(string format)
var ContextKeyIPAddress = ContextKey("ipaddr")

//ContextKeyIP context key for reqeust ip
var ContextKeyIP = ContextKey("ip")
