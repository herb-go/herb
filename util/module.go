package util

import (
	"fmt"
	"sort"
)

type Module struct {
	Name    string
	Handler func()
}

var PrintModuleInitLog = false

type modulelist []Module

func (m modulelist) Len() int {
	return len(m)
}
func (m modulelist) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
func (m modulelist) Less(i, j int) bool {
	return m[i].Name < m[j].Name
}

var Modules = modulelist{}

func RegisteModule(name string, handler func()) Module {
	m := Module{Name: name, Handler: handler}
	Modules = append(Modules, m)
	return m
}

func InitModulesOrderByName() {
	sort.Sort(Modules)
	for k := range Modules {
		if PrintModuleInitLog {
			fmt.Println("Herb-go util debug: Init module " + Modules[k].Name)
		}
		Modules[k].Handler()
	}
}
