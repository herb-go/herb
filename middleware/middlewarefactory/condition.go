package middlewarefactory

import "net/http"

type Condition interface {
	CheckRequest(*http.Request) (bool, error)
}

func Not(r *http.Request, c Condition) (bool, error) {
	ok, err := c.CheckRequest(r)
	if err != nil {
		return false, err
	}
	return !ok, nil
}

func And(r *http.Request, c ...Condition) (bool, error) {
	for k := range c {
		ok, err := c[k].CheckRequest(r)
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
		ok, err := c[k].CheckRequest(r)
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
func (c *PlainCondition) CheckRequest(r *http.Request) (bool, error) {
	var result bool
	var err error
	if c.Disabled {
		return false, nil
	}

	if len(c.Conditions) != 0 {
		conditions := append([]Condition{c.Condition}, c.Conditions...)
		if c.Or {
			result, err = Or(r, conditions...)
		} else {
			result, err = And(r, conditions...)
		}
	} else {
		result, err = c.Condition.CheckRequest(r)
	}
	if err != nil {
		return false, err
	}
	if c.Not {
		result = !result
	}
	return result, nil
}
