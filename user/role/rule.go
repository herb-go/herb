package role

//Rule user authorize rule interface
type Rule interface {
	Execute(roles ...*Role) (bool, error)
}

//RuleOr rules commbined in or operation.
type RuleOr struct {
	Rules []Rule
}

//Execute execute as rule provider.
//Execute will success if any rule success
func (c *RuleOr) Execute(roles ...*Role) (bool, error) {
	for _, v := range c.Rules {
		result, err := v.Execute(roles...)
		if err != nil {
			return false, err
		}
		if result {
			return true, nil
		}
	}
	return false, nil
}

//RuleAnd rules commbined in not operation.
type RuleAnd struct {
	Rules []Rule
}

//Execute will success if all rule success
func (c *RuleAnd) Execute(roles ...*Role) (bool, error) {
	for _, v := range c.Rules {
		result, err := v.Execute(roles...)
		if err != nil {
			return false, err
		}
		if !result {
			return false, nil
		}
	}
	return true, nil
}

//RuleNot a not operation to rule
type RuleNot struct {
	Rule Rule
}

//Execute will success if rule fail
func (c *RuleNot) Execute(roles ...*Role) (bool, error) {
	result, err := c.Rule.Execute(roles...)
	if err != nil {
		return false, err
	}
	return !result, nil
}

//Not Create a not rule
func Not(c Rule) *RuleNot {
	return &RuleNot{
		Rule: c,
	}
}

//And create a and rule
func And(c ...Rule) *RuleAnd {
	return &RuleAnd{
		Rules: c,
	}
}

//Or create a or rule
func Or(c ...Rule) *RuleOr {
	return &RuleOr{
		Rules: c,
	}
}

//RuleSet a set of rule
type RuleSet struct {
	Rule Rule
}

//NewRuleSet create new rule set
func NewRuleSet(Rule Rule) *RuleSet {
	return &RuleSet{
		Rule: Rule,
	}
}

//Not set ruleset to not operated rule.
func (ruleset *RuleSet) Not() *RuleSet {
	ruleset.Rule = Not(ruleset.Rule)
	return ruleset
}

//And combine ruleset with rules by and operate
func (ruleset *RuleSet) And(c ...Rule) *RuleSet {
	rs := make([]Rule, len(c)+1)
	rs[0] = ruleset.Rule
	copy(rs[1:], c)
	ruleset.Rule = And(rs...)
	return ruleset
}

//Or combine ruleset with rules by or operate
func (ruleset *RuleSet) Or(c ...Rule) *RuleSet {
	rs := make([]Rule, len(c)+1)
	rs[0] = ruleset.Rule
	copy(rs[1:], c)
	ruleset.Rule = Or(rs...)
	return ruleset
}

//Execute exccute ruleset
func (ruleset *RuleSet) Execute(roles ...*Role) (bool, error) {
	return ruleset.Rule.Execute(roles...)
}
