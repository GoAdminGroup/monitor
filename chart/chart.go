package chart

import (
	"github.com/GoAdminGroup/monitor/chart/data"
	"github.com/GoAdminGroup/monitor/dashboard"
	"github.com/GoAdminGroup/monitor/dashboard/param"
)

type BaseChart struct {
	DataSource  dashboard.DataSource
	Type        dashboard.ChartType
	Name        string
	Title       string
	Description string
	Param       param.Param
	DataFormat  *data.Format
	Size        []int
	Height      int
	Id          int
}

func (b *BaseChart) GetDataSource() dashboard.DataSource {
	return b.DataSource
}

func (b *BaseChart) GetSize() []int {
	return b.Size
}

func (b *BaseChart) GetHeight() int {
	return b.Height
}

func (b *BaseChart) SetSize(value []int) {
	b.Size = value
}

func (b *BaseChart) GetId() int {
	return b.Id
}

func (b *BaseChart) SetId(id int) {
	b.Id = id
}

func (b *BaseChart) SetHeight(value int) {
	b.Height = value
}

type BaseConfig struct {
	Title       string
	Description string
	Param       param.Param
	Size        []int
	Height      int
	DataFormat  *data.Format
}
