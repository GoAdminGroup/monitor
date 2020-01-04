package chart

import (
	"fmt"
	"github.com/GoAdminGroup/components/echarts"
	"github.com/GoAdminGroup/go-admin/modules/config"
	template2 "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/monitor/dashboard"
	"github.com/GoAdminGroup/monitor/dashboard/param"
	template3 "github.com/GoAdminGroup/monitor/template"
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

	graphData, err := b.getData(param.Combine(b.Param))

	if err != nil {
		return "", err
	}

	line := b.getLine(graphData)

	return template2.Get(config.Get().Theme).Box().
		SetHeader(template.HTML("<b>" + b.Title + "</b>" + fmt.Sprintf(template3.Get(config.Get().Theme).GetGraphBtn(), b.Id))).
		SetBody(echarts.NewChart().SetContent(line).GetContent()).
		SetHeadColor("#f8f9fb").
		GetContent(), nil
}

func (b *Graph) getLine(graphData *dashboard.GraphData) *charts.Line {

	if len(graphData.YAxisList) > 0 && len(graphData.YAxisList[0]) > 0 {
		b.DataFormat.CorrectUnit(graphData.YAxisList[0][0])
	} else {
		b.DataFormat.CorrectUnit(0)
	}

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

	yopts := charts.YAxisOpts{
		AxisLabel: charts.LabelTextOpts{Formatter: "{value} " + b.DataFormat.Flag()},
		SplitLine: charts.SplitLineOpts{
			Show: true,
			LineStyle: charts.LineStyleOpts{
				Type: "dashed",
			},
		},
	}
	xopts := charts.XAxisOpts{
		//AxisLabel: charts.LabelTextOpts{Formatter: "{value} " + b.DataFormat.Type.Flag(b.DataFormat.Unit)},
		SplitLine: charts.SplitLineOpts{
			Show: true,
			LineStyle: charts.LineStyleOpts{
				Type: "dashed",
			},
		},
	}

	if b.YMax == -1 || b.YMin == -1 {
		if b.DataFormat.Type.IsPercent() {
			yopts.Max = 100
			yopts.Min = 0
		}
	} else {
		yopts.Max = b.YMax
		yopts.Min = b.YMin
	}

	line := charts.NewLine()
	line.SetGlobalOptions(yopts, xopts,
		charts.ToolboxOpts{Show: true},
		charts.LegendOpts{Left: "1%", Top: "2%"},
		charts.TooltipOpts{Show: true, Trigger: "axis"}, //  Formatter: "{c} " + b.DataFormat.Flag()
	)
	line.SetSeriesOptions(
		charts.LineOpts{Smooth: true},
	)
	line.AddXAxis(labels)

	for i := 0; i < len(graphData.YAxisList); i++ {
		if len(graphData.YAxisList[i]) > 0 {
			line.AddYAxis(yLabels[i], b.DataFormat.FormatDatas(graphData.YAxisList[i]),
				//charts.ColorOpts(color_scheme.DefaultScheme),
			)
		}
	}

	line.Width = "100%"
	line.Height = fmt.Sprintf("%dpx", b.Height)
	return line
}

func (b *Graph) getData(param param.Param) (*dashboard.GraphData, error) {
	data, err := b.DataSource.GetData(param)
	if err != nil {
		return nil, err
	}

	return dashboard.GraphDataFromChartData(data), nil
}

func (b *Graph) GetData(param param.Param) (template.JS, error) {

	graphData, err := b.getData(param.Combine(b.Param))

	if err != nil {
		return "", err
	}

	line := b.getLine(graphData)

	return echarts.NewChart().SetContent(line).GetOptions(), nil
}
