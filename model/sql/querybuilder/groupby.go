package querybuilder

import "strings"

// GroupByClause group by clause struct
type GroupByClause struct {
	Buidler *Builder
	Fields  []string
}

// Add add fields to group by clause
func (q *GroupByClause) Add(fields ...string) *GroupByClause {
	q.Fields = append(q.Fields, fields...)
	return q
}

// QueryCommand return query command
func (q *GroupByClause) QueryCommand() string {
	if len(q.Fields) == 0 {
		return ""
	}
	return "GROUP BY " + strings.Join(q.Fields, " , ")
}

// QueryArgs return query args.
func (q *GroupByClause) QueryArgs() []interface{} {
	return []interface{}{}
}

//NewGroupByClause create  new group by clause.
func (b *Builder) NewGroupByClause() *GroupByClause {
	return &GroupByClause{
		Buidler: b,
		Fields:  []string{},
	}
}
