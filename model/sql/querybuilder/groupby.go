package querybuilder

import "strings"

type GroupByClause struct {
	buidler *Builder
	Fields  []string
}

func (q *GroupByClause) Add(fields ...string) *GroupByClause {
	q.Fields = append(q.Fields, fields...)
	return q
}
func (q *GroupByClause) QueryCommand() string {
	if len(q.Fields) == 0 {
		return ""
	}
	return " " + strings.Join(q.Fields, " , ") + " "
}
func (q *GroupByClause) QueryArgs() []interface{} {
	return []interface{}{}
}

func (b *Builder) NewGroupByClause() *GroupByClause {
	return &GroupByClause{
		buidler: b,
		Fields:  []string{},
	}
}
