package httpinfo

type Controller interface {
	Write([]byte)
	BeforeWriteHeader()
	BeforeWrite()
	Error() error
}

type NopController struct {
}

func (p *NopController) Write([]byte) {
}

func (p *NopController) BeforeWriteHeader() {
}

func (p *NopController) BeforeWrite() {
}

func (p *NopController) Error() error {
	return nil
}

var DefaultController Controller = &NopController{}
