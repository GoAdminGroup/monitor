package template

type Template interface {
	GetDashboardStyle() string
	GetGraphBtn() string
	GetSingleStatBtn() string
	GetSingleStatContent() string
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