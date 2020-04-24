package dashboard

import (
	"fmt"
	"html/template"
	"strconv"

	"github.com/GoAdminGroup/go-admin/modules/config"
	template2 "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/monitor/dashboard/param"
	template3 "github.com/GoAdminGroup/monitor/template"
)

type Dashboard interface {
	GetContent(params []param.Param) (template.HTML, error)
	GetCharts() ChartList
	SetKey(key string) Dashboard
	GetKey() string
	GetChart(id int) Chart
	GetTitle() template.HTML
	GetDescription() template.HTML
}

type Layout uint8

type Gen func() Dashboard

type GeneratorList map[string]Dashboard

var list = make(GeneratorList)

func Add(name string, board Dashboard) {
	if _, duplicate := list[name]; duplicate {
		panic("has been registered")
	}
	list[name] = board.SetKey(name)
}

func Get(name string) Dashboard {
	if value, ok := list[name]; ok {
		return value
	}
	panic("wrong dashboard name")
}

type DefaultDashboard struct {
	Layout      Layout
	Charts      ChartList
	DS          DSList
	Title       template.HTML
	Description template.HTML
	Key         string
}

func NewDefaultDashboard(title, desc template.HTML) *DefaultDashboard {
	return &DefaultDashboard{Title: title, Description: desc}
}

func (d *DefaultDashboard) SetKey(key string) Dashboard {
	d.Key = key
	return d
}

func (d *DefaultDashboard) GetKey() string {
	return d.Key
}

func (d *DefaultDashboard) GetChart(id int) Chart {
	return d.Charts[id]
}

func (d *DefaultDashboard) GetContent(params []param.Param) (template.HTML, error) {

	var (
		width      = 0
		height     = 0
		content    template.HTML
		row        template.HTML
		components = template2.Get(config.Get().Theme)
		col        template.HTML
	)

	for i, ga := range d.Charts {

		p := make(param.Param, 0)
		if len(params) > 1 {
			p = params[i+1]
		}

		if ga.GetSize()[1] == 12 {

			if height != 0 {
				if col != template.HTML("") {
					row += components.Col().SetSize(map[string]string{"md": strconv.Itoa(d.Charts[i-1].GetSize()[0])}).SetContent(col).GetContent()
					col = ""
				}
				width += d.Charts[i-1].GetSize()[0]
			}

			height = 0

			// newline
			if width+ga.GetSize()[0] > 12 {
				width = 0
				if row != template.HTML("") {
					content += components.Row().SetContent(row).GetContent()
					row = ""
				}
			} else {
				width += ga.GetSize()[0]
			}

			c, err := ga.GetContent(SetChartType(p, ga.GetType()))
			if err == nil {
				row += components.Col().
					SetSize(map[string]string{"md": strconv.Itoa(ga.GetSize()[0])}).
					SetContent(c).GetContent()
			}
		} else {

			if height+ga.GetSize()[1] > 12 {
				height = ga.GetSize()[1]
				if col != template.HTML("") {
					row += components.Col().SetSize(map[string]string{"md": strconv.Itoa(d.Charts[i-1].GetSize()[0])}).SetContent(col).GetContent()
					col = ""
				}
				width += d.Charts[i-1].GetSize()[0]
			} else {
				height += ga.GetSize()[1]
			}

			c, err := ga.GetContent(SetChartType(p, ga.GetType()))
			if err == nil {
				col += components.Row().SetContent(c).GetContent()
			}

			if i == len(d.Charts)-1 {
				row += components.Col().SetSize(map[string]string{"md": strconv.Itoa(ga.GetSize()[0])}).SetContent(col).GetContent()
			}
		}

		if i == len(d.Charts)-1 {
			content += components.Row().SetContent(row).GetContent()
		}
	}

	tmpl := template3.Get(config.Get().Theme)

	return tmpl.GetToolBar(params[0]["interval"].(string)) + content +
		template.HTML(fmt.Sprintf(template3.Get(config.Get().Theme).GetDashboardStyle(), config.Get().Prefix(), d.Key)), nil
}

func (d *DefaultDashboard) GetTitle() template.HTML {
	return d.Title
}

func (d *DefaultDashboard) GetDescription() template.HTML {
	return d.Description
}

func (d *DefaultDashboard) GetCharts() ChartList {
	return d.Charts
}

func (d *DefaultDashboard) AddChart(g Chart) {
	g.SetId(len(d.Charts))
	d.Charts = d.Charts.Add(g)
}

func (d *DefaultDashboard) AddChartWithDSKey(dskey string, g Chart) {
	g.SetId(len(d.Charts))
	d.Charts = d.Charts.Add(g.SetDataSource(d.DS.FindByName(dskey)))
}

func (d *DefaultDashboard) AddDataSource(key string, param param.Param) {
	so := GetDataSource(key)

	if err := so.Init(param); err != nil {
		fmt.Println("add data source fail: ", err)
	}

	d.DS = d.DS.Add(so)
}

type DataSource interface {
	GetData(param param.Param) (ChartData, error)
	Init(param param.Param) error
	GetName() string
}

type DSList []DataSource

func (l DSList) FindByName(name string) DataSource {
	for _, ds := range l {
		if ds.GetName() == name {
			return ds
		}
	}
	panic("wrong datasource name")
}

func (l DSList) Add(source DataSource) DSList {
	return append(l, source)
}

type DSmap map[string]DataSource

func (d DSmap) Get(key string) DataSource {
	return d[key]
}

func GetDataSource(key string) DataSource {
	return dscenter.Get(key)
}

var dscenter = make(DSmap)

func RegisterDS(key string, source DataSource) {
	if _, duplicate := dscenter[key]; duplicate {
		panic("has been registered")
	}
	dscenter[key] = source
}

type ChartList []Chart

func (list ChartList) Add(g Chart) ChartList {
	return append(list, g)
}

type ChartType uint8

const (
	Graph ChartType = iota
	Gauge
	SingleStat
)

func (g ChartType) IsGraph() bool {
	return g == Graph
}

func (g ChartType) IsSingleStat() bool {
	return g == SingleStat || g == Gauge
}

const ParamChartTypeKey = "gauge_type"

func NewWithChartType(t ChartType) param.Param {
	var p = make(param.Param)
	p[ParamChartTypeKey] = t
	return p
}

func SetChartType(p param.Param, t ChartType) param.Param {
	p[ParamChartTypeKey] = t
	return p
}

func GetChartType(p param.Param) ChartType {
	return p[ParamChartTypeKey].(ChartType)
}

type Chart interface {
	SetDataSource(DataSource) Chart
	GetDataSource() DataSource
	GetType() ChartType
	GetSize() []int
	GetId() int
	SetId(id int)
	GetHeight() int
	SetSize([]int)
	SetHeight(int)
	GetData(param param.Param) (template.JS, error)
	GetContent(param param.Param) (template.HTML, error)
}

type ChartData interface{}

func GraphDataFromChartData(g ChartData) *GraphData {
	return g.(*GraphData)
}

func SingleStatDataFromChartData(g ChartData) SingleStatData {
	return g.(SingleStatData)
}

type GraphData struct {
	XAxis      []string
	YAxisList  [][]float64
	yaxisIndex int
}

func NewGraphData() *GraphData {
	return &GraphData{
		XAxis:     make([]string, 0),
		YAxisList: make([][]float64, 0),
	}
}

func (g *GraphData) AddLabels(label string) {
	g.XAxis = append(g.XAxis, label)
}

func (g *GraphData) AddYAxisPoint(index int, yaxis float64) {
	g.YAxisList[index] = append(g.YAxisList[index], yaxis)
}

func (g *GraphData) AddYAxis() {
	g.YAxisList = append(g.YAxisList, make([]float64, 0))
}

type SingleStatData float64
