package chart

import (
	"fmt"
	"github.com/GoAdminGroup/components/echarts"
	"github.com/GoAdminGroup/go-admin/modules/config"
	template2 "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/monitor/dashboard"
	"github.com/GoAdminGroup/monitor/dashboard/param"
	"github.com/go-echarts/go-echarts/charts"
	"html/template"
)

type Graph struct {
	BaseChart
	XLabels []string
	YLabels []string
	YMax    float64
	YMin    float64
}

type GraphConfig struct {
	BaseConfig
	XLabels []string
	YLabels []string
	YMax    float64
	YMin    float64
}

func NewGraph(cfg GraphConfig) *Graph {

	ymax := cfg.YMax
	ymin := cfg.YMin

	if ymax == 0 {
		ymax = -1
	}
	if ymin == 0 {
		ymin = -1
	}

	return &Graph{
		BaseChart: BaseChart{
			Title:       cfg.Title,
			Description: cfg.Description,
			Size:        cfg.Size,
			Height:      cfg.Height,
			Param:       cfg.Param,
			DataFormat:  cfg.DataFormat,
		},
		XLabels: cfg.XLabels,
		YLabels: cfg.YLabels,
		YMax:    ymax,
		YMin:    ymin,
	}
}

func (b *Graph) SetDataSource(ds dashboard.DataSource) dashboard.Chart {
	b.DataSource = ds
	return b
}

func (b *Graph) GetType() dashboard.ChartType {
	return dashboard.Graph
}

func (b *Graph) GetContent(param param.Param) (template.HTML, error) {

	param = param.Combine(b.Param)

	data, err := b.DataSource.GetData(param)
	if err != nil {
		return "", err
	}

	graphData := dashboard.GraphDataFromChartData(data)

	labels := b.XLabels
	if len(labels) == 0 {
		labels = graphData.XAxis
	}
	yLabels := b.YLabels
	if len(yLabels) == 0 {
		for i := 0; i < len(graphData.YAxisList); i++ {
			yLabels = append(yLabels, fmt.Sprintf("%dxxxxx", i))
		}
	}

	opts := charts.YAxisOpts{
		AxisLabel: charts.LabelTextOpts{Formatter: "{value} " + b.DataFormat.Type.Flag()},
		SplitLine: charts.SplitLineOpts{
			Show: true,
			LineStyle: charts.LineStyleOpts{
				Type: "dashed",
			},
		},
	}

	if b.YMax == -1 || b.YMin == -1 {
		if b.DataFormat.Type.IsPercent() {
			opts.Max = 100
			opts.Min = 0
		}
	} else {
		opts.Max = b.YMax
		opts.Min = b.YMin
	}

	line := charts.NewLine()
	line.SetGlobalOptions(opts, charts.ToolboxOpts{Show: true}, charts.LegendOpts{Left: "1%", Top: "2%"})
	line.AddXAxis(labels)

	for i := 0; i < len(graphData.YAxisList); i++ {
		line.AddYAxis(yLabels[i], formatDatas(graphData.YAxisList[i], b.DataFormat), charts.LabelTextOpts{Show: true})
	}

	line.Width = "100%"
	line.Height = fmt.Sprintf("%dpx", b.Height)

	b.Id = line.ChartID

	components := template2.Get(config.Get().Theme)

	return components.Box().
		SetHeader(template.HTML("<b>" + b.Title + "</b>" + btn)).
		SetBody(echarts.NewChart().SetContent(line).GetContent()).
		SetHeadColor("#f8f9fb").
		GetContent(), nil
}
