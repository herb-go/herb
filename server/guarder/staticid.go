package guarder

var DefaultStaticID = "staticid"

type StaticID string

func (i StaticID) ID() string {
	if i == "" {
		return DefaultStaticID
	}
	return string(i)
}

func (i StaticID) IsEmpty() bool {
	return i == ""
}

func (i StaticID) Equal(v string) bool {
	return string(i) == v
}
