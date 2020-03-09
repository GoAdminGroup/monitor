package template

import "html/template"

type Template interface {
	GetDashboardStyle() string
	GetGraphBtn() string
	GetSingleStatBtn() string
	GetSingleStatContent() string
	GetToolBar(string) template.HTML
}

var list = map[string]Template{
	"adminlte": new(Adminlte),
}

func Get(key string) Template {
	if v, ok := list[key]; ok {
		return v
	}
	panic("wrong key")
}
