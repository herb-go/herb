package ui

//Component basic ui component
type Component interface {
	ComponentID() string
}

//ComponentLabels ui component with labels
type ComponentLabels interface {
	GetComponentLabels() Labels
	SetComponentLabels(Labels)
	GetLabel(string) string
}

// LabelsComponent component struct with labels
type LabelsComponent struct {
	labels Labels
}

//ComponentID return component id
func (c *LabelsComponent) ComponentID() string {
	return ""
}

// GetComponentLabels get labels from component
func (c *LabelsComponent) GetComponentLabels() Labels {
	return c.labels
}

//SetComponentLabels get labels from component
func (c *LabelsComponent) SetComponentLabels(labels Labels) {
	c.labels = labels
}

//GetLabel get label from componet labels
func (c *LabelsComponent) GetLabel(field string) string {
	if c.labels == nil {
		return ""
	}
	l := c.labels.GetLabel(field)
	if l == "" {
		return ""
	}
	return l
}
