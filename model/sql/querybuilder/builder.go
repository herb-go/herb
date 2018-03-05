package builder

type Builder struct {
	Driver string
}

var DefaultBuilder = &Builder{
	Driver: "",
}
