package guarder

type ParamsKey string

func (k ParamsKey) LoadFrom(p *Params) string {
	return (*p)[string(k)]
}

func (k ParamsKey) SetTo(p *Params, v string) {
	(*p)[string(k)] = v
}

const ParamsKeyID = ParamsKey("id")
const ParamsKeyCredential = ParamsKey("token")

type Params map[string]string

func (p *Params) ID() string {
	return ParamsKeyID.LoadFrom(p)
}

func (p *Params) SetID(v string) {
	ParamsKeyID.SetTo(p, v)
}

func (p *Params) Credential() string {
	return ParamsKeyCredential.LoadFrom(p)
}

func (p *Params) SetCredential(v string) {
	ParamsKeyCredential.SetTo(p, v)
}

func NewParams() *Params {
	return &Params{}
}
