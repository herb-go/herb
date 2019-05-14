package querybuilder

import (
	"strings"
)

type order struct {
	Field string
	Asc   bool
}
type OrderByClause struct {
	buidler *Builder
	orders  []order
}

func (q *OrderByClause) Add(field string, asc bool) *OrderByClause {
	q.orders = append(q.orders, order{
		Field: field,
		Asc:   asc,
	})
	return q
}
func (q *OrderByClause) QueryCommand() string {
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
func (q *OrderByClause) QueryArgs() []interface{} {
	return []interface{}{}
}

func (b *Builder) NewOrderByClause() *OrderByClause {
	return &OrderByClause{
		buidler: b,
		orders:  []order{},
	}
}
