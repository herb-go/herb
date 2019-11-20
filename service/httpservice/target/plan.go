package target

import (
	"net/http"
)

type PlainPlan struct {
	Doer
	Target
	*Builders
}

func (p *PlainPlan) Execute() (*http.Response, error) {
	return Execute(p)
}
func (p *PlainPlan) WithDoder(d Doer) *PlainPlan {
	p.Doer = d
	return p
}

func (p *PlainPlan) WithTarget(t Target) *PlainPlan {
	p.Target = t
	return p
}
func (p *PlainPlan) WithBuilders(b *Builders) *PlainPlan {
	p.Builders = b
	return p
}

func NewPlan() *PlainPlan {
	return &PlainPlan{
		Builders: NewBuilders(),
	}
}

type Plan interface {
	Doer
	Target
	Builder
}

func Execute(p Plan) (*http.Response, error) {
	return Do(p, p, p)
}
