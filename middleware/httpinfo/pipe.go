package httpinfo

type Pipe interface {
	Write([]byte)
	Check()
	Discarded() bool
	Error() error
}

type NopPipe struct {
}

func (p *NopPipe) Write([]byte) {
}

func (p *NopPipe) Check() {
}

func (p *NopPipe) Discarded() bool {
	return true
}

func (p *NopPipe) Error() error {
	return nil
}

var DefaultPipe Pipe = &NopPipe{}
