package middlewarefactory

import (
	"net/http"
	"time"
)

type TimeCondition struct {
	Start int64
	End   int64
}

func (c *TimeCondition) CheckRequest(*http.Request) (bool, error) {
	now := time.Now().Unix()
	if (c.Start <= 0 || now >= c.Start) && (c.End <= 0 || now <= c.End) {
		return true, nil
	}
	return false, nil
}

var TimeConditionFactory = ConditionFactoryFunc(func(loader func(v interface{}) error) (Condition, error) {
	c := &TimeCondition{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	return c, nil
})
