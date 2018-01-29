package role

type Rule interface {
	Execute(roles ...Role) (bool, error)
}

type RuleOr struct {
	Rules []Rule
}

func (c *RuleOr) Execute(roles ...Role) (bool, error) {
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

type RuleAnd struct {
	Rules []Rule
}

func (c *RuleAnd) Execute(roles ...Role) (bool, error) {
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

type RuleNot struct {
	Rule Rule
}

func (c *RuleNot) Execute(roles ...Role) (bool, error) {
	result, err := c.Rule.Execute(roles...)
	if err != nil {
		return false, err
	}
	return !result, nil
}

func Not(c Rule) *RuleNot {
	return &RuleNot{
		Rule: c,
	}
}

func And(c ...Rule) *RuleAnd {
	return &RuleAnd{
		Rules: c,
	}
}

func Or(c ...Rule) *RuleOr {
	return &RuleOr{
		Rules: c,
	}
}

type RuleSet struct {
	Rule Rule
}

func NewRuleSet(Rule Rule) *RuleSet {
	return &RuleSet{
		Rule: Rule,
	}
}
func (c *RuleSet) Not() *RuleSet {
	c.Rule = Not(c.Rule)
	return c
}

func (ruleset *RuleSet) And(c ...Rule) *RuleSet {
	rs := make([]Rule, len(c)+1)
	rs[0] = ruleset.Rule
	copy(rs[1:], c)
	ruleset.Rule = And(rs...)
	return ruleset
}

func (ruleset *RuleSet) Or(c ...Rule) *RuleSet {
	rs := make([]Rule, len(c)+1)
	rs[0] = ruleset.Rule
	copy(rs[1:], c)
	ruleset.Rule = Or(rs...)
	return ruleset
}

func (c *RuleSet) Execute(roles ...Role) (bool, error) {
	return c.Rule.Execute(roles...)
}
