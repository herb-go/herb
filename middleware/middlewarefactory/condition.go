package middlewarefactory

import "net/http"

type Condition interface {
	MatchRequest(*http.Request) (bool, error)
}
type ConditionFunc func(*http.Request) (bool, error)

func (f ConditionFunc) MatchRequest(r *http.Request) (bool, error) {
	return f(r)
}

var EmptyCondition = ConditionFunc(func(*http.Request) (bool, error) {
	return true, nil
})

type ConditionFactory func(func(v interface{}) error) (Condition, error)

func Not(r *http.Request, c Condition) (bool, error) {
	ok, err := c.MatchRequest(r)
	if err != nil {
		return false, err
	}
	return !ok, nil
}

func And(r *http.Request, c ...Condition) (bool, error) {
	for k := range c {
		ok, err := c[k].MatchRequest(r)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

func Or(r *http.Request, c ...Condition) (bool, error) {
	for k := range c {
		ok, err := c[k].MatchRequest(r)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

type PlainCondition struct {
	Condition  Condition
	Not        bool
	Or         bool
	Disabled   bool
	Conditions []Condition
}

func NewPlainCondition() *PlainCondition {
	return &PlainCondition{
		Conditions: []Condition{},
	}
}
func (c *PlainCondition) MatchRequest(r *http.Request) (bool, error) {
	var result bool
	var err error
	if c.Disabled {
		return false, nil
	}
	if c.Condition == nil {
		c.Condition = EmptyCondition
	}
	if len(c.Conditions) != 0 {
		conditions := append([]Condition{c.Condition}, c.Conditions...)
		if c.Or {
			result, err = Or(r, conditions...)
		} else {
			result, err = And(r, conditions...)
		}
	} else {
		result, err = c.Condition.MatchRequest(r)
	}
	if err != nil {
		return false, err
	}
	if c.Not {
		result = !result
	}
	return result, nil
}
