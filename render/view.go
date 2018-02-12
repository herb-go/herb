package render

import "net/http"

//NamedView named view main struct.
type NamedView struct {
	//Name view name.
	Name string
	//Renderer view renderer.
	Renderer *Renderer
}

//Render render view with given data to response.
//Return bytes length wrote and any error if raised.
func (v *NamedView) Render(writer http.ResponseWriter, data interface{}) (int, error) {
	return v.RenderError(writer, data, http.StatusOK)
}

//MustRender render view with given data to response.
//Return bytes length wrote.
//Panic if any error raised.
func (v *NamedView) MustRender(writer http.ResponseWriter, data interface{}) int {
	return v.MustRenderError(writer, data, http.StatusOK)
}

//RenderError render view with given data and status code to response.
//Return bytes length wrote and any error if raised.
func (v *NamedView) RenderError(writer http.ResponseWriter, data interface{}, status int) (int, error) {
	output, err := v.RenderBytes(data)
	if err != nil {
		return 0, err
	}
	return WriteHTML(writer, []byte(output), status)
}

//MustRenderError render view with given data and status code to response.
//Return bytes length wrote.
//Panic if any error raised.
func (v *NamedView) MustRenderError(writer http.ResponseWriter, data interface{}, status int) int {
	output := v.MustRenderBytes(data)
	return MustWriteHTML(writer, []byte(output), status)
}

//RenderBytes render view with given data to bytes.
//Return bytes length wrote and any error if raised.
func (v *NamedView) RenderBytes(data interface{}) ([]byte, error) {
	return v.Renderer.Execute(v.Name, data)
}

//MustRenderBytes render view with given data to bytes.
//Return bytes length wrote.
//Panic if any error raised.
func (v *NamedView) MustRenderBytes(data interface{}) []byte {
	output, err := v.RenderBytes(data)
	if err != nil {
		panic(err)
	}
	return output
}
