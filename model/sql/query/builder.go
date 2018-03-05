package query

type Builder struct {
	Driver string
}

var DefaultBuilder = &Builder{
	Driver: "",
}
