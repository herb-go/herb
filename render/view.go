package render

import "net/http"

type NamedView struct {
	Name     string
	Renderer *Renderer
}

func (v *NamedView) Render(writer http.ResponseWriter, data interface{}) (int, error) {
	return v.RenderError(writer, data, http.StatusOK)
}
func (v *NamedView) MustRender(writer http.ResponseWriter, data interface{}) int {
	return v.MustRenderError(writer, data, http.StatusOK)
}
func (v *NamedView) RenderError(writer http.ResponseWriter, data interface{}, status int) (int, error) {
	output, err := v.RenderString(data)
	if err != nil {
		return 0, err
	}
	return WriteHTML(writer, []byte(output), status)
}
func (v *NamedView) MustRenderError(writer http.ResponseWriter, data interface{}, status int) int {
	output := v.MustRenderString(data)
	return MustWriteHTML(writer, []byte(output), status)
}
func (v *NamedView) RenderString(data interface{}) (string, error) {
	var cv = v.Renderer.view(v.Name)
	if cv == nil {
		return "", NewViewError(v.Name, ErrorViewNotExist)
	}
	return cv.Execute(data)
}

func (v *NamedView) MustRenderString(data interface{}) string {
	output, err := v.RenderString(data)
	if err != nil {
		panic(err)
	}
	return output
}
