package role

import "testing"

func TestRole(t *testing.T) {
	var rolename = "editor"
	var rolename2 = "admin"
	var rulenameNotOwned1 = "no1"
	// var rolenameNotOwned2 = "no2"
	var fieldname1 = "country"
	var roleDataUS = "US"
	var roleDataCN = "CN"
	var fieldname2 = "field2"
	var roleData2 = "data2"
	var role = New(rolename)
	var role2 = New(rolename2)
	var ruleNotOwned1 = New(rulenameNotOwned1)
	var roleWithData2 = New(rolename)
	roleWithData2.AddData(fieldname1, roleDataUS, roleDataCN)
	roleWithData2.AddData(fieldname2, roleData2)
	var rule = New(rolename)
	result, _ := rule.Execute()
	if result != false {
		t.Error(result)
	}
	result, _ = rule.Execute(*role)
	if result != true {
		t.Error(result)
	}
	result, _ = rule.Execute(*role, *role2)
	if result != true {
		t.Error(result)
	}
	result, _ = ruleNotOwned1.Execute(*role, *role2)
	if result != false {
		t.Error(result)
	}
	var roleCN = New(rolename)
	roleCN.AddData(fieldname1, roleDataCN)
	var ruleCN = New(rolename)
	ruleCN.AddData(fieldname1, roleDataCN)
	result, _ = ruleCN.Execute(*roleCN)
	if result != true {
		t.Error(result)
	}
	var ruleUS = New(rolename)
	ruleUS.AddData(fieldname1, roleDataUS)
	result, _ = ruleUS.Execute(*roleCN)
	if result != false {
		t.Error(result)
	}
	result, _ = ruleUS.Execute(*role)
	if result != false {
		t.Error(result)
	}
	var ruleCNUS = New(rolename)
	ruleCNUS.AddData(fieldname1, roleDataUS, roleDataCN)
	result, _ = ruleCNUS.Execute(*roleCN)
	if result != false {
		t.Error(result)
	}
	var roleUSCN = New(rolename)
	roleUSCN.AddData(fieldname1, roleDataCN, roleDataUS)
	result, _ = ruleCNUS.Execute(*roleUSCN)
	if result != true {
		t.Error(result)
	}
	result, _ = ruleCN.Execute(*roleUSCN)
	if result != true {
		t.Error(result)
	}
	result, _ = ruleUS.Execute(*roleUSCN)
	if result != true {
		t.Error(result)
	}
	result, _ = ruleUS.Execute(*roleWithData2)
	if result != true {
		t.Error(result)
	}
	result, _ = roleWithData2.Execute(*ruleCN)
	if result != false {
		t.Error(result)
	}
}

func TestRules(t *testing.T) {
	var result bool
	var rolename = "editor"
	// var rolenameNotOwned2 = "no2"
	var fieldname1 = "country"
	var roleDataUS = "US"
	var roleDataCN = "CN"
	var fieldname2 = "field2"
	var roleData2 = "data2"
	var rolename2 = "admin"
	var role = New(rolename)
	var role2 = New(rolename2)
	var roleCN = New(rolename)
	roleCN.AddData(fieldname1, roleDataCN)
	var ruleCN = New(rolename)
	ruleCN.AddData(fieldname1, roleDataCN)
	var roleUS = New(rolename)
	roleUS.AddData(fieldname1, roleDataUS)
	var ruleUS = New(rolename)
	ruleUS.AddData(fieldname1, roleDataUS)
	var roleUSCN = New(rolename)
	roleUSCN.AddData(fieldname1, roleDataCN, roleDataUS)
	var ruleCNUS = New(rolename)
	ruleCNUS.AddData(fieldname1, roleDataUS, roleDataCN)
	var rulesCN = Roles{*ruleCN}
	var rulesCNUS = Roles{*ruleCN, *ruleUS}
	var roleWithData2 = New(rolename)
	roleWithData2.AddData(fieldname1, roleDataUS, roleDataCN)
	roleWithData2.AddData(fieldname2, roleData2)
	var ruleWithData2 = Roles{*roleWithData2}
	result, _ = NewRoles().Execute()
	if result != true {
		t.Error(result)
	}
	result, _ = rulesCN.Execute()
	if result != false {
		t.Error(result)
	}
	var rules = Roles{*role}
	result, _ = rules.Execute(*role)
	if result != true {
		t.Error(result)
	}
	result, _ = rules.Execute(*role2)
	if result != false {
		t.Error(result)
	}
	result, _ = rulesCN.Execute(*roleCN)
	if result != true {
		t.Error(result)
	}
	result, _ = rulesCN.Execute(*roleUS, *roleCN)
	if result != true {
		t.Error(result)
	}
	result, _ = rulesCN.Execute(*roleUS)
	if result != false {
		t.Error(result)
	}
	result, _ = rulesCN.Execute(*roleWithData2)
	if result != true {
		t.Error(result)
	}
	result, _ = rulesCNUS.Execute(*roleCN)
	if result != true {
		t.Error(result)
	}
	result, _ = rulesCNUS.Execute(*roleUS, *roleCN)
	if result != true {
		t.Error(result)
	}
	result, _ = rulesCNUS.Execute(*roleUS)
	if result != true {
		t.Error(result)
	}
	result, _ = ruleWithData2.Execute(*roleUSCN)
	if result != false {
		t.Error(result)
	}
	result, _ = ruleWithData2.Execute(*roleUS, *roleCN)
	if result != false {
		t.Error(result)
	}
	result, _ = ruleWithData2.Execute(*roleWithData2)
	if result != true {
		t.Error(result)
	}
}
