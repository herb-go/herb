package builder

type Fields map[string]interface{}

func (f Fields) Set(field string, v interface{}) Fields {
	f[field] = v
	return f
}
