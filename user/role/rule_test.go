package role

import "testing"

type fixedRule struct {
	result bool
}

func (r *fixedRule) Execute(roles ...Role) (bool, error) {
	return r.result, nil
}

func TestRule(t *testing.T) {
	var trueRule = &fixedRule{result: true}
	var falseRule = &fixedRule{result: false}
	var ruleNotTrue = Not(trueRule)
	var result, _ = ruleNotTrue.Execute()
	if result != false {
		t.Error(result)
	}
	var ruleNotFalse = Not(falseRule)
	result, _ = ruleNotFalse.Execute()
	if result != true {
		t.Error(result)
	}
	var ruleAndAllTrue = And(trueRule, trueRule, trueRule, trueRule)
	result, _ = ruleAndAllTrue.Execute()
	if result != true {
		t.Error(result)
	}
	var ruleAndWithFalse = And(trueRule, trueRule, trueRule, falseRule)
	result, _ = ruleAndWithFalse.Execute()
	if result != false {
		t.Error(result)
	}
	var ruleOrAllFalse = Or(falseRule, falseRule, falseRule, falseRule)
	result, _ = ruleOrAllFalse.Execute()
	if result != false {
		t.Error(result)
	}
	var ruleOrWithTrue = Or(trueRule, falseRule, falseRule, falseRule)
	result, _ = ruleOrWithTrue.Execute()
	if result != true {
		t.Error(result)
	}
}

func TestRuleSet(t *testing.T) {
	var trueRule = &fixedRule{result: true}
	var falseRule = &fixedRule{result: false}
	var ruleSetNotTrue = NewRuleSet(trueRule).Not()
	var result, _ = ruleSetNotTrue.Execute()
	if result != false {
		t.Error(result)
	}
	var ruleSetNotFalse = NewRuleSet(falseRule).Not()
	result, _ = ruleSetNotFalse.Execute()
	if result != true {
		t.Error(result)
	}
	var ruleSetAndAllTrue = NewRuleSet(trueRule).And(trueRule, trueRule, trueRule)
	result, _ = ruleSetAndAllTrue.Execute()
	if result != true {
		t.Error(result)
	}
	var ruleSetAndWithFalse = NewRuleSet(trueRule).And(trueRule, trueRule, falseRule)
	result, _ = ruleSetAndWithFalse.Execute()
	if result != false {
		t.Error(result)
	}
	var ruleSetOrAllFalse = NewRuleSet(falseRule).Or(falseRule, falseRule, falseRule)
	result, _ = ruleSetOrAllFalse.Execute()
	if result != false {
		t.Error(result)
	}
	var ruleSetOrWithTrue = NewRuleSet(trueRule).Or(falseRule, falseRule, falseRule)
	result, _ = ruleSetOrWithTrue.Execute()
	if result != true {
		t.Error(result)
	}
	var ruleSetChain = NewRuleSet(trueRule).
		And(
			trueRule, trueRule, trueRule,
			Or(
				trueRule,
				falseRule,
			),
		).
		Or(
			And(
				trueRule,
				falseRule,
			),
		)
	result, _ = ruleSetChain.Execute()
	if result != true {
		t.Error(result)
	}
}
