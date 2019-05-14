package querybuilder

type LimitClause struct {
	buidler *Builder
	limit   *int
	offset  *int
}

func (q *LimitClause) SetOffset(o int) *LimitClause {
	offset := o
	q.offset = &offset
	return q
}
func (q *LimitClause) SetLimit(l int) *LimitClause {
	limit := l
	q.limit = &limit
	return q
}

func (q *LimitClause) QueryCommand() string {

	if q.limit == nil && q.offset == nil {
		return ""
	}
	return q.buidler.LoadDriver().LimitCommandBuilder(q)
}
func (q *LimitClause) QueryArgs() []interface{} {
	if q.limit == nil && q.offset == nil {
		return nil
	}
	return q.buidler.LoadDriver().LimitArgBuilder(q)
}

func (b *Builder) NewLimitClause() *LimitClause {
	return &LimitClause{
		buidler: b,
	}
}
