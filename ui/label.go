package ui

// Label label interface
type Label interface {
	Label() string
}

//Labels labels interface
type Labels interface {
	GetLabel(string) string
}

//MapLabels kabels collection in map form
type MapLabels map[string]string

//GetLabel get field label by field name
func (l MapLabels) GetLabel(field string) string {
	label, ok := l[field]
	if ok == false {
		return ""
	}
	return label

}

// StringLabel string label
type StringLabel string

//Label label string
func (l StringLabel) Label() string {
	return string(l)
}
