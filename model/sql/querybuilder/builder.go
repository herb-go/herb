package querybuilder

import (
	"fmt"
	"sync"
)

type Builder struct {
	Driver string
	driver BuilderDriver
	lock   sync.Mutex
}

func NewBuilder() *Builder {
	return &Builder{}
}
func (b *Builder) LoadDriver() BuilderDriver {
	d := b.driver
	if d != nil {
		return d
	}
	b.lock.Lock()
	defer b.lock.Unlock()
	d = loadBuilderDriver(b.Driver)
	b.driver = d
	return b.driver
}

var DefaultBuilder = &EmptyBuilderDriver{}

type BuilderDriver interface {
	LimitQueryBuilder(q *LimitQuery) string
	LimitArgBuilder(q *LimitQuery) []interface{}
}

type EmptyBuilderDriver struct {
}

func (d *EmptyBuilderDriver) LimitQueryBuilder(q *LimitQuery) string {
	var command = ""
	if q.limit != nil {
		command = "LIMIT ? "
	}
	if q.offset != nil {
		command += "OFFSET ? "
	}
	return command
}

func (d *EmptyBuilderDriver) LimitArgBuilder(q *LimitQuery) []interface{} {
	var args = []interface{}{}
	if q.limit != nil {
		args = append(args, q.limit)
	}
	if q.offset != nil {
		args = append(args, q.offset)
	}
	return args
}

var drivers = map[string]BuilderDriver{}

func RegisterBuilder(name string, driver BuilderDriver) {
	drivers[name] = driver
	if Debug && driver != nil {
		fmt.Println("querybuilder: build driver '" + name + "' registered.")
	}
}

func loadBuilderDriver(name string) BuilderDriver {
	d := drivers[name]
	if d == nil {
		return DefaultBuilder
	}
	return d
}

func init() {
	RegisterBuilder("mysql", DefaultBuilder)
}
