package db

type Optimizer interface {
	MustOptimize(cmd string, args []interface{}) (string, []interface{})
}

type OptimizerFactory func(func(v interface{}) error) (Optimizer, error)

type DefaultOptimizerConfig struct {
	Replace *ReplaceOptimizer
}

var DefaultOptimizerFactory = func(loader func(v interface{}) error) (Optimizer, error) {
	if loader == nil {
		return nil, nil
	}
	c := &DefaultOptimizerConfig{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	if c.Replace == nil {
		return nil, err
	}
	return c.Replace, nil
}

type ReplaceOptimizer map[string]string

func (o *ReplaceOptimizer) MustOptimize(cmd string, args []interface{}) (string, []interface{}) {
	data, ok := (*o)[cmd]
	if ok {
		cmd = data
	}
	return cmd, args
}
