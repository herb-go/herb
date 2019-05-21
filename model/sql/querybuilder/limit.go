package querybuilder

// LimitClause limit clause struct
type LimitClause struct {
	buidler *Builder
	// Limit sql limit arg
	Limit *int
	// Offset sql offset arg
	Offset *int
}

// SetOffset set limit clause offset
func (q *LimitClause) SetOffset(o int) *LimitClause {
	offset := o
	q.Offset = &offset
	return q
}

//SetLimit set  limit clause linut
func (q *LimitClause) SetLimit(l int) *LimitClause {
	limit := l
	q.Limit = &limit
	return q
}

// QueryCommand return query command
func (q *LimitClause) QueryCommand() string {

	if q.Limit == nil && q.Offset == nil {
		return ""
	}
	return q.buidler.LoadDriver().LimitCommandBuilder(q)
}

// QueryArgs return query args
func (q *LimitClause) QueryArgs() []interface{} {
	if q.Limit == nil && q.Offset == nil {
		return nil
	}
	return q.buidler.LoadDriver().LimitArgBuilder(q)
}

// NewLimitClause create new limit clause
func (b *Builder) NewLimitClause() *LimitClause {
	return &LimitClause{
		buidler: b,
	}
}
