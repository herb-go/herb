package requestform

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/herb-go/herb/model"
)

const MsgBadRequest = "Bad request."

func MustValidateJSONRequest(r *http.Request, m RequestValidator) bool {
	var body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if len(body) > 0 {
		err = json.Unmarshal(body, &m)
		if err != nil {
			m.SetBadRequest(true)
			m.AddErrorf("", MsgBadRequest)
			return false
		}
	}
	err = m.InitWithRequest(r)
	if err != nil {
		panic(err)
	}
	return model.MustValidate(m)

}
func MustRenderErrors(w http.ResponseWriter, m RequestValidator) {
	if m.BadRequest() {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(422)
	bytes, err := json.Marshal(m.Errors())
	if err != nil {
		panic(err)
	}
	_, err = w.Write(bytes)
	if err != nil {
		panic(err)
	}
}

type Form struct {
	model.Model
	badRequest bool
}

func (f *Form) InitWithRequest(*http.Request) error {
	return nil
}

func (model *Form) BadRequest() bool {
	return model.badRequest
}
func (model *Form) SetBadRequest(v bool) {
	model.badRequest = v
}

func (model *Form) HasError() bool {
	return model.badRequest || len(model.Errors()) != 0
}

type RequestValidator interface {
	model.Validator
	InitWithRequest(*http.Request) error
	BadRequest() bool
	SetBadRequest(v bool)
}
