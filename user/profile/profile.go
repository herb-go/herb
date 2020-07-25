package profile

type Profile []*Field

func (p *Profile) Data() []*Field {
	return []*Field(*p)
}

func (p *Profile) Clone() *Profile {
	result := make(Profile, len(*p))
	copy(result, *p)
	return &result
}
func (p *Profile) With(name string, value string) *Profile {
	f := NewField()
	f.Name = name
	f.Value = value
	return p.WithFields(f)
}
func (p *Profile) WithFields(f ...*Field) *Profile {
	*p = append(*p, f...)
	return p
}

func (p *Profile) Load(name string) string {
	length := len(*p)
	if length == 0 {
		return ""
	}
	for i := length - 1; i >= 0; i-- {
		if (*p)[i].Name == name {
			return (*p)[i].Value
		}
	}
	return ""
}

func (p *Profile) LoadAny(names ...string) string {
	for k := range names {
		v := p.Load(names[k])
		if v != "" {
			return v
		}
	}
	return ""
}

func (p *Profile) Filter(names ...string) *Profile {
	result := Profile{}
Fields:
	for _, field := range *p {
		for _, name := range names {
			if field.Name == name {
				result = append(result, field)
				continue Fields
			}
		}
	}
	return &result
}

func (p *Profile) Chain(profiles ...*Profile) *Profile {
	for k := range profiles {
		p.WithFields(profiles[k].Data()...)
	}
	return p
}
func (p *Profile) LoadAllByName(name string) []string {
	result := []string{}
	for _, v := range *p {
		if v.Name == name {
			result = append(result, v.Value)
		}
	}
	return result
}

// NewProfile create mew profile
func NewProfile() *Profile {
	return &Profile{}
}
