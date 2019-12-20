package chart

import (
	"fmt"
	"github.com/GoAdminGroup/monitor/dashboard"
	"github.com/GoAdminGroup/monitor/dashboard/param"
	"strconv"
)

type BaseChart struct {
	DataSource  dashboard.DataSource
	Type        dashboard.ChartType
	Name        string
	Title       string
	Description string
	Param       param.Param
	DataFormat  Format
	Size        []int
	Height      int
	Id          string
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

func (b *BaseChart) GetId() string {
	return b.Id
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
	DataFormat  Format
}

type Format struct {
	Type     DataType
	Unit     int
	Decimals int
}

type DataType uint8

const (
	Currency DataType = iota
	Percent
	Number
	Second
)

func (d DataType) IsPercent() bool {
	return d == Percent
}

func formatDatas(data []float64, format Format) []float64 {
	if format.Type == Percent {
		for i := 0; i < len(data); i++ {
			data[i] = decimal(data[i]*100, format.Decimals)
		}
	} else {
		if format.Decimals != 0 {
			for i := 0; i < len(data); i++ {
				data[i] = decimal(data[i], format.Decimals)
			}
		}
	}
	return data
}

func formatData(data float64, format Format) string {
	if format.Type == Second {
		if data > 60 && data < 3600 {
			return decimalString(data/60, format.Decimals) + "分钟"
		} else if data > 3600 && data < 86400 {
			return decimalString(data/3600, format.Decimals) + "小时"
		} else if data > 86400 {
			return decimalString(data/86400, format.Decimals) + "天"
		}
	}
	return decimalString(data, format.Decimals)
}

func (d DataType) Flag() string {
	switch d {
	case Percent:
		return "%"
	}
	return ""
}

func decimal(value float64, dec int) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%."+strconv.Itoa(dec)+"f", value), 64)
	return value
}

func decimalString(value float64, dec int) string {
	return fmt.Sprintf("%."+strconv.Itoa(dec)+"f", value)
}

const btn = `<div class="zoom-in-btn" style="position: relative;cursor: pointer;color: #a7a7a7;
float: right;right: 0.5%;"><i class="fa fa-arrows-alt"></i></div><div class="zoom-out-btn" style="display:none;position: relative;cursor: pointer;color: #a7a7a7;
float: right;right: 0.5%;"><i class="fa fa-compress"></i></div>`
