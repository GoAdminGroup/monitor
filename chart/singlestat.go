package chart

import (
	"fmt"
	"github.com/GoAdminGroup/go-admin/modules/config"
	template2 "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/monitor/dashboard"
	"github.com/GoAdminGroup/monitor/dashboard/param"
	"html/template"
)

type SingleStat struct {
	BaseChart
	Color string
}

type SingleStatConfig struct {
	BaseConfig
	Color string
}

func NewSingleStat(cfg SingleStatConfig) *SingleStat {
	return &SingleStat{
		BaseChart: BaseChart{
			Title:       cfg.Title,
			Description: cfg.Description,
			Size:        cfg.Size,
			Height:      cfg.Height,
			Param:       cfg.Param,
			DataFormat:  cfg.DataFormat,
		},
		Color: cfg.Color,
	}
}

func (s *SingleStat) SetDataSource(ds dashboard.DataSource) dashboard.Chart {
	s.DataSource = ds
	return s
}

func (s *SingleStat) GetType() dashboard.ChartType {
	return dashboard.SingleStat
}

func (s *SingleStat) GetContent(param param.Param) (template.HTML, error) {

	param = param.Combine(s.Param)

	data, err := s.DataSource.GetData(param)
	if err != nil {
		return "", err
	}

	singleStatData := dashboard.SingleStatDataFromChartData(data)

	components := template2.Get(config.Get().Theme)

	content := `<div style="height:%dpx;line-height: %dpx;font-size: 2.5em;font-weight:bold;text-align:center;color: %s;">%s</div>`

	return components.Box().
		SetHeader(template.HTML("<b>" + s.Title + "</b>" + btn)).
		SetBody(template.HTML(fmt.Sprintf(content, s.Height, s.Height, s.Color, formatData(float64(singleStatData), s.DataFormat)))).
		SetHeadColor("#f8f9fb").
		GetContent(), nil
}
