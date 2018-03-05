package querybuilder

func (b *Builder) concatWith(separator string, q ...Query) *PlainQuery {
	var query = b.New("")
	for k := range q {
		if q[k] == nil {
			continue
		}
		command := q[k].QueryCommand()
		if command != "" {
			query.Command += q[k].QueryCommand() + separator
		}
		query.Args = append(query.Args, q[k].QueryArgs()...)
	}
	if query.Command != "" {
		query.Command = query.Command[:len(query.Command)-len(separator)]
	}
	return query
}
func (b *Builder) Concat(q ...Query) *PlainQuery {
	return b.concatWith(" ", q...)
}

func (b *Builder) Comma(q ...Query) *PlainQuery {
	return b.concatWith(" , ", q...)
}
func (b *Builder) Lines(q ...Query) *PlainQuery {
	return b.concatWith("\n", q...)
}
func (b *Builder) And(q ...Query) *PlainQuery {
	if (len(q)) == 1 {
		return b.New(q[0].QueryCommand(), q[0].QueryArgs()...)
	}
	var query = b.concatWith(" AND ", q...)
	if query.Command != "" {
		query.Command = "( " + query.Command + " )"
	}
	return query
}

func (b *Builder) Or(q ...Query) *PlainQuery {
	if (len(q)) == 1 {
		return b.New(q[0].QueryCommand(), q[0].QueryArgs()...)
	}
	var query = b.concatWith(" OR ", q...)
	if query.Command != "" {
		query.Command = "( " + query.Command + " )"
	}
	return query
}
