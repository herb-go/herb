package role

import "net/http"

//Roles type role list
type Roles []*Role

func getValueMapCache(roles Roles) map[int]map[string]map[string]bool {
	var valuemap = map[int]map[string]map[string]bool{}
	for roleindex, role := range roles {
		if valuemap[roleindex] == nil {
			valuemap[roleindex] = map[string]map[string]bool{}
		}
		for fieldname, fielddata := range role.Data {
			if valuemap[roleindex][fieldname] == nil {
				valuemap[roleindex][fieldname] = map[string]bool{}
			}
			for _, value := range fielddata {
				valuemap[roleindex][fieldname][value] = true
			}
		}
	}
	return valuemap
}

//Rule implement rule provdier interface.
func (rules *Roles) Rule(*http.Request) (Rule, error) {
	return rules, nil
}

//Add add role to roles
func (rules *Roles) Add(r *Role) *Roles {
	*rules = append(*rules, r)
	return rules
}

//Execute execute role as rule provider.
//If rule data if empty,any role with same name will success.
//Otherwise,only roles with data which covers all rule data will success.
//Method will suceess if All rules success.
func (rules *Roles) Execute(roles ...*Role) (bool, error) {
	if len(*rules) == 0 {
		return true, nil
	}
	if len(roles) == 0 {
		return false, nil
	}
	valuemap := getValueMapCache(roles)
	for _, rule := range *rules {
	NextRole:
		for roleindex, role := range roles {
			if rule.Name == role.Name {
				if rule.Data == nil {
					return true, nil
				}
				if rule.Data != nil {
					for fieldname := range rule.Data {
						for _, ruledata := range rule.Data[fieldname] {
							if valuemap[roleindex][fieldname][ruledata] == false {
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
			//Name not matched,check next role
		}
		//All role match fail,check next Rule
	}
	//All Match fail
	return false, nil
}

//NewRoles create new roles with given rolesnames.
//You can't use this method to create roles with data.
func NewRoles(rolenames ...string) *Roles {
	var roles = make(Roles, len(rolenames))
	for k := range rolenames {
		roles[k] = New(rolenames[k])
	}
	return &roles
}
