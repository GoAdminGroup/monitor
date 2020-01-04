package chart

import (
	"fmt"
	"github.com/GoAdminGroup/go-admin/modules/config"
	template2 "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/monitor/dashboard"
	"github.com/GoAdminGroup/monitor/dashboard/param"
	template3 "github.com/GoAdminGroup/monitor/template"
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
		return s.content("N/A"), nil
	}

	singleStatData := dashboard.SingleStatDataFromChartData(data)

	s.DataFormat.CorrectUnit(float64(singleStatData))

	return s.content(s.DataFormat.FormatData(float64(singleStatData))), nil
}

func (s *SingleStat) content(data string) template.HTML {
	return template2.Get(config.Get().Theme).Box().
		SetHeader(template.HTML("<b>" + s.Title + "</b>" + fmt.Sprintf(template3.Get(config.Get().Theme).GetSingleStatBtn(), s.Id))).
		SetBody(template.HTML(fmt.Sprintf(template3.Get(config.Get().Theme).GetSingleStatContent(), s.Height, s.Height, s.Color, data))).
		SetHeadColor("#f8f9fb").
		GetContent()
}

func (s *SingleStat) GetData(param param.Param) (template.JS, error) {
	return "", nil
}
