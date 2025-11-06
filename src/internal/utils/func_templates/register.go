package functemplates

import (
	"html/template"
)

var funcMap template.FuncMap

func register(name string, fn any) {
	funcMap[name] = fn
}

func init() {
	funcMap = make(template.FuncMap)
	register("default", defaultFunc)
}



func GetFuncMap() template.FuncMap {

	return funcMap
}


