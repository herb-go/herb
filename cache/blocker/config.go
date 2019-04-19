package blocker

import (
	"time"
)

//Rule blocker block rule
type Rule struct {
	StatusCode       int
	Limit            int64
	DurationInSecond int64
}

//ApplyTo apply block rule to blocker
func (r *Rule) ApplyTo(b *Blocker) error {
	b.Block(r.StatusCode, r.Limit, time.Duration(r.DurationInSecond)*time.Second)
	return nil
}

//Rules blocker block rule list
type Rules []*Rule

//ApplyTo  apply rule list to blocker
func (r *Rules) ApplyTo(b *Blocker) error {
	for _, v := range *r {
		err := v.ApplyTo(b)
		if err != nil {
			return err
		}
	}
	return nil
}

//NewRules create new rule list
func NewRules() *Rules {
	var r Rules = []*Rule{}
	return &r
}
