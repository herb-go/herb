package role

//Role user role main struct
type Role struct {
	//Name role name
	Name string
	//Data role data
	Data map[string][]string
}

//AddData add data to role's special field
func (r *Role) AddData(field string, data ...string) {
	if r.Data == nil {
		r.Data = map[string][]string{}
	}
	if r.Data[field] == nil {
		r.Data[field] = []string{}
	}
	r.Data[field] = append(r.Data[field], data...)
}

//New create new role with role name
func New(name string) *Role {
	return &Role{
		Name: name,
		Data: nil,
	}
}

//Execute execute role as rule provider.
//If rule data if empty,any role with same name will success.
//Otherwise,only roles with data which covers all rule data will success.
func (r *Role) Execute(roles ...Role) (bool, error) {
	if len(roles) == 0 {
		return false, nil
	}
NextRole:
	for _, role := range roles {
		if r.Name == role.Name {
			if r.Data == nil {
				return true, nil
			}
			if r.Data != nil {
				for fieldname := range r.Data {
					var valuemap = map[string]bool{}
					for _, value := range role.Data[fieldname] {
						valuemap[value] = true
					}
					for _, ruledata := range r.Data[fieldname] {
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
