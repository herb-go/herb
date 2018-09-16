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
	var command = ""
	if q.limit == nil && q.offset == nil {
		return command
	}
	switch q.buidler.Driver {
	default:
		if q.limit != nil {
			command = "LIMIT ? "
		}
		if q.offset != nil {
			command += "OFFSET ? "
		}
	}
	return command
}
func (q *LimitQuery) QueryArgs() []interface{} {
	args := []interface{}{}
	if q.limit == nil && q.offset == nil {
		return args
	}
	switch q.buidler.Driver {
	default:
		if q.limit != nil {
			args = append(args, q.limit)
		}
		if q.offset != nil {
			args = append(args, q.offset)
		}
	}
	return args
}

func (b *Builder) NewLimitQuery() *LimitQuery {
	return &LimitQuery{
		buidler: b,
	}
}
