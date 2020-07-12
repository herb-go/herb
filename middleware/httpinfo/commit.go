package httpinfo

import "net/http"

type CommitController struct {
	request   *http.Request
	response  *Response
	lasterror error
	checker   Validator
	NopController
}

func (c *CommitController) BeforeWriteHeader() {
	c.check()
}

func (c *CommitController) BeforeWrite() {
	c.check()
}

func (c *CommitController) check() {
	if c.response.autocommit != false {
		return
	}
	ok, err := c.checker.Validate(c.request, c.response)
	if err != nil {
		c.lasterror = err
		c.response.Commit()
		return
	}
	if !ok {
		c.response.Commit()
	}
}
func (c *CommitController) WithChecker(v Validator) *CommitController {
	c.checker = v
	return c
}

func NewCommitController(req *http.Request, resp *Response) *CommitController {
	return &CommitController{
		checker:  ValidatorAlways,
		request:  req,
		response: resp,
	}
}
