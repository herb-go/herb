package querybuilder

type LimitQuery struct {
	buidler *Builder
	limit   *int
	offset  *int
}

func (q *LimitQuery) SetOffset(o int) *LimitQuery {
	offset := o
	q.offset = &offset
	return q
}
func (q *LimitQuery) SetLimit(l int) *LimitQuery {
	limit := l
	q.limit = &limit
	return q
}

func (q *LimitQuery) QueryCommand() string {

	if q.limit == nil && q.offset == nil {
		return ""
	}
	return q.buidler.LoadDriver().LimitCommandBuilder(q)
}
func (q *LimitQuery) QueryArgs() []interface{} {
	if q.limit == nil && q.offset == nil {
		return nil
	}
	return q.buidler.LoadDriver().LimitArgBuilder(q)
}

func (b *Builder) NewLimitQuery() *LimitQuery {
	return &LimitQuery{
		buidler: b,
	}
}
