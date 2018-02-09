package query

func concatWith(separator string, q ...Query) *PlainQuery {
	var query = New("")
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
func Concat(q ...Query) *PlainQuery {
	return concatWith(" ", q...)
}

func Comma(q ...Query) *PlainQuery {
	return concatWith(" , ", q...)
}
func Lines(q ...Query) *PlainQuery {
	return concatWith("\n", q...)
}
func And(q ...Query) *PlainQuery {
	if (len(q)) == 1 {
		return New(q[0].QueryCommand(), q[0].QueryArgs()...)
	}
	var query = concatWith(" AND ", q...)
	if query.Command != "" {
		query.Command = "( " + query.Command + " )"
	}
	return query
}

func Or(q ...Query) *PlainQuery {
	if (len(q)) == 1 {
		return New(q[0].QueryCommand(), q[0].QueryArgs()...)
	}
	var query = concatWith(" OR ", q...)
	if query.Command != "" {
		query.Command = "( " + query.Command + " )"
	}
	return query
}
