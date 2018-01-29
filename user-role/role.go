package role

type Role struct {
	Name string
	Data map[string][]string
}

func (r *Role) AddData(field string, data ...string) {
	if r.Data == nil {
		r.Data = map[string][]string{}
	}
	if r.Data[field] == nil {
		r.Data[field] = []string{}
	}
	r.Data[field] = append(r.Data[field], data...)
}

func New(name string) *Role {
	return &Role{
		Name: name,
		Data: nil,
	}
}

func (rule *Role) Execute(roles ...Role) (bool, error) {
	if len(roles) == 0 {
		return false, nil
	}
NextRole:
	for _, role := range roles {
		if rule.Name == role.Name {
			if rule.Data == nil {
				return true, nil
			}
			if rule.Data != nil {
				for fieldname := range rule.Data {
					var valuemap = map[string]bool{}
					for _, value := range role.Data[fieldname] {
						valuemap[value] = true
					}
					for _, ruledata := range rule.Data[fieldname] {
						if valuemap[ruledata] == false {
							//Data not matched.
							//Field matched ,check next role
							continue NextRole
						}
					}
					//All rule data in field matched .
					//Field matched.
				}
				//All field matched.Role matched
				return true, nil
			}
		}
	}
	return false, nil
}
