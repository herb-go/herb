package db

type Optimizer interface {
	Optimize(cmd string, args []interface{}) (string, []interface{}, error)
}
