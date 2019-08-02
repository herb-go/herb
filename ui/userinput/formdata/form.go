package formdata

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/herb-go/herb/ui/userinput"
)

//MsgBadRequest field error msg used in unmarshaler request form fail.
const MsgBadRequest = "Bad request."

//MustValidateRequestBody unmarshal form with request body with given Unmarshaler, then init form with request and  validate it.
//Return validate result.
//Add bad request error to form empty field if unmarshal form fail.
//Panci if any other error raised.
func MustValidateRequestBody(r *http.Request, Unmarshaler func([]byte, interface{}) error, m RequestValidator) bool {
	var body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if len(body) > 0 {
		err = Unmarshaler(body, &m)
		if err != nil {
			m.SetBadRequest(true)
			m.AddErrorf("", MsgBadRequest)
			return false
		}
	}
	m.SetHTTPRequest(r)
	err = m.InitWithRequest(r)
	if err != nil {
		panic(err)
	}
	return model.MustValidate(m)
}

//MustValidateJSONRequest unmarshal form with request body with json.unmarshal, then init form with request and  validate it.
//Return validate result.
//Add bad request error to form empty field if unmarshal form fail.
//Panci if any other error raised.
func MustValidateJSONRequest(r *http.Request, m RequestValidator) bool {
	return MustValidateRequestBody(r, json.Unmarshal, m)
}

//MustRenderError render error data of model to http response.
//Data was marshaled by marshaler.
//Http bad request error will be executed if model has a bad request error.
//Return bytes length rendered
//Panic if any error raised.
func MustRenderError(w http.ResponseWriter, marshaler func(interface{}) ([]byte, error), model RequestValidator) int {
	if model.BadRequest() {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return 0
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(422)
	bytes, err := marshaler(model.Errors())
	if err != nil {
		panic(err)
	}
	result, err := w.Write(bytes)
	if err != nil {
		panic(err)
	}
	return result

}

//MustRenderErrorsJSON render error data of model to http response.
//Data was marshaled by json.marshal.
//Http bad request error will be executed if model has a bad request error.
//Return bytes length rendered
//Panic if any error raised.
func MustRenderErrorsJSON(w http.ResponseWriter, model RequestValidator) {
	MustRenderError(w, json.Marshal, model)
}

//Form form struct
type Form struct {
	model.Model
	badRequest  bool
	httprequest *http.Request
}

//InitWithRequest init model with given request.
//Return any error if rasied.
//You can override this methon in your own form.
func (model *Form) InitWithRequest(*http.Request) error {
	return nil
}

//BadRequest return if form has a bad request error.
func (model *Form) BadRequest() bool {
	return model.badRequest
}

//SetBadRequest set whether form has a bad request error.
func (model *Form) SetBadRequest(hasError bool) {
	model.badRequest = hasError
}

//SetHTTPRequest Set http request to form
func (model *Form) SetHTTPRequest(r *http.Request) {
	model.httprequest = r
}

//HTTPRequest Return http request in form
func (model *Form) HTTPRequest() *http.Request {
	return model.httprequest
}

//HasError return if model has any error.
//Return true if form has a bad request error.
func (model *Form) HasError() bool {
	return model.badRequest || model.Model.HasError()
}

//RequestValidator interface of request form that can be validated.
type RequestValidator interface {
	model.Validator
	//InitWithRequest init model with given request.
	//Return any error if rasied.
	InitWithRequest(*http.Request) error
	//BadRequest return if form has a bad request error.
	BadRequest() bool
	//SetBadRequest set whether form has a bad request error.
	SetBadRequest(v bool)
	//Return http request in form
	HTTPRequest() *http.Request
	//Set http request to form
	SetHTTPRequest(*http.Request)
}
