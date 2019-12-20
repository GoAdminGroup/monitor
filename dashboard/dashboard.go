package dashboard

import (
	"fmt"
	"github.com/GoAdminGroup/go-admin/modules/config"
	template2 "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/monitor/dashboard/param"
	"html/template"
	"strconv"
)

type Dashboard interface {
	GetLayout() Layout
	GetContent(params []param.Param) (template.HTML, error)
	GetCharts() ChartList
	GetTitle() string
	GetDescription() string
}

type Layout uint8

type Gen func() Dashboard

type GeneratorList map[string]Dashboard

var list = make(GeneratorList)

func Add(name string, board Dashboard) {
	if _, duplicate := list[name]; duplicate {
		panic("has been registered")
	}
	list[name] = board
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
	Title       string
	Description string
}

func NewDefaultDashboard(title, desc string) *DefaultDashboard {
	return &DefaultDashboard{Title: title, Description: desc}
}

func (d *DefaultDashboard) GetLayout() Layout {
	return d.Layout
}

func (d *DefaultDashboard) GetContent(params []param.Param) (template.HTML, error) {

	var (
		width      = 0
		height     = 0
		content    template.HTML
		row        template.HTML
		components = template2.Get(config.Get().Theme)
		chartIds   = make([]string, 0)
		col        template.HTML
	)

	for i, ga := range d.Charts {

		p := make(param.Param, 0)
		if len(params) > 0 {
			p = params[i]
		}

		fmt.Println("ga.GetSize()", ga.GetSize())
		fmt.Println("width", width)
		fmt.Println("height", height)

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
				chartIds = append(chartIds, ga.GetId())
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
				chartIds = append(chartIds, ga.GetId())
			}

			if i == len(d.Charts)-1 {
				row += components.Col().SetSize(map[string]string{"md": strconv.Itoa(ga.GetSize()[0])}).SetContent(col).GetContent()
			}
		}

		if i == len(d.Charts)-1 {
			content += components.Row().SetContent(row).GetContent()
		}
	}

	return content + sideBarControlJS, nil
}

const sideBarControlJS = `<script>
	$("body").addClass("sidebar-collapse")
	$(".zoom-in-btn").on('click', function (event) {
		$(this).parent().parent().parent().addClass("zoom-in-container")
		let container = $(this).parent().next().children().children()
		container.attr("data-raw-height", container.height())
		container.css("height", "680px");
		let chartID = "myChart_" + container.attr("id")
		eval(chartID + ".resize()")
		$(this).hide()
		$(this).next().show()
	});
	$(".zoom-out-btn").on('click', function (event) {
		$(this).parent().parent().parent().removeClass("zoom-in-container")
		let container = $(this).parent().next().children().children() 
		container.css("height", container.attr("data-raw-height"));
		let chartID = "myChart_" + container.attr("id")
		eval(chartID + ".resize()")
		$(this).hide()
		$(this).prev().show()
	});
</script>
<style>
	.echarts-container{margin-top:0px;}
	.zoom-in-container{position:absolute;width:93.5%;height:82%;z-index:999;}
	.row {margin-right: 0;margin-left: 0;}
	.col-lg-1, .col-lg-10, .col-lg-11, .col-lg-12, .col-lg-2, .col-lg-3, .col-lg-4, .col-lg-5, .col-lg-6, .col-lg-7, .col-lg-8, .col-lg-9, .col-md-1, .col-md-10, .col-md-11, .col-md-12, .col-md-2, .col-md-3, .col-md-4, .col-md-5, .col-md-6, .col-md-7, .col-md-8, .col-md-9, .col-sm-1, .col-sm-10, .col-sm-11, .col-sm-12, .col-sm-2, .col-sm-3, .col-sm-4, .col-sm-5, .col-sm-6, .col-sm-7, .col-sm-8, .col-sm-9, .col-xs-1, .col-xs-10, .col-xs-11, .col-xs-12, .col-xs-2, .col-xs-3, .col-xs-4, .col-xs-5, .col-xs-6, .col-xs-7, .col-xs-8, .col-xs-9{padding-right: 6px;padding-left: 6px;}
	.box {margin-bottom: 6px;}
</style>`

func (d *DefaultDashboard) GetTitle() string {
	return d.Title
}

func (d *DefaultDashboard) GetDescription() string {
	return d.Description
}

func (d *DefaultDashboard) GetCharts() ChartList {
	return d.Charts
}

func (d *DefaultDashboard) AddChart(g Chart) {
	d.Charts = d.Charts.Add(g)
}

func (d *DefaultDashboard) AddChartWithDSKey(dskey string, g Chart) {
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
	GetId() string
	GetHeight() int
	SetSize([]int)
	SetHeight(int)
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
