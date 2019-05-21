package querybuilder

import (
	"strings"
)

// Order order data struct
type Order struct {
	// Field order field
	Field string
	//  Asc order field in asc or not
	Asc bool
}

// OrderByClause order by clause struct
type OrderByClause struct {
	buidler *Builder
	// Orders order list
	Orders []Order
}

// Add add new order to clause with given field and order.
func (q *OrderByClause) Add(field string, asc bool) *OrderByClause {
	q.Orders = append(q.Orders, Order{
		Field: field,
		Asc:   asc,
	})
	return q
}

// QueryCommand return query command
func (q *OrderByClause) QueryCommand() string {
	if len(q.Orders) == 0 {
		return ""
	}
	commands := make([]string, len(q.Orders))
	for k, v := range q.Orders {
		commands[k] = v.Field
		if v.Asc {
			commands[k] += " ASC "
		} else {
			commands[k] += " DESC "
		}
	}
	return "ORDER BY " + strings.Join(commands, " , ")
}

// QueryArgs return query args
func (q *OrderByClause) QueryArgs() []interface{} {
	return []interface{}{}
}

// NewOrderByClause create order by clause
func (b *Builder) NewOrderByClause() *OrderByClause {
	return &OrderByClause{
		buidler: b,
		Orders:  []Order{},
	}
}
