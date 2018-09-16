package querybuilder

import (
	"strings"
)

type order struct {
	Field string
	Asc   bool
}
type OrderByQuery struct {
	buidler *Builder
	orders  []order
}

func (q *OrderByQuery) Add(field string, asc bool) *OrderByQuery {
	q.orders = append(q.orders, order{
		Field: field,
		Asc:   asc,
	})
	return q
}
func (q *OrderByQuery) QueryCommand() string {
	if len(q.orders) == 0 {
		return ""
	}
	commands := make([]string, len(q.orders))
	for k, v := range q.orders {
		commands[k] = v.Field
		if v.Asc {
			commands[k] += " ASC "
		} else {
			commands[k] += " DESC "
		}
	}
	return "ORDER BY " + strings.Join(commands, " , ")
}
func (q *OrderByQuery) QueryArgs() []interface{} {
	return []interface{}{}
}

func (b *Builder) NewOrderByQuery() *OrderByQuery {
	return &OrderByQuery{
		buidler: b,
		orders:  []order{},
	}
}
