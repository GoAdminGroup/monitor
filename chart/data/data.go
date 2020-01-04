package data

import (
	"fmt"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/monitor/chart/data/unit"
	"github.com/GoAdminGroup/monitor/constant"
	"strconv"
)

type Format struct {
	Type     Type
	DSUnit   unit.Unit
	Unit     unit.Unit
	Decimals int
}

func (f *Format) FormatDatas(data []float64) []float64 {
	for i := 0; i < len(data); i++ {
		data[i] = decimal(data[i]*f.DSUnit.Coefficient()/f.Unit.Coefficient(), f.Decimals)
	}
	return data
}

func (f *Format) FormatData(data float64) string {
	return decimalString(data*f.DSUnit.Coefficient()/f.Unit.Coefficient(), f.Decimals) +
		language.GetWithScope(f.Unit.String(), constant.Scope)
}

func (f *Format) CorrectUnit(data float64) {
	if !f.DSUnit.IsValid() {
		f.DSUnit = f.GetDefaultUnit(data)
	}
	if !f.Unit.IsValid() {
		f.Unit = unit.GetUnitFromData(data, f.DSUnit)
	}
}

func (f *Format) Flag() string {
	return f.Type.Flag(f.Unit)
}

func (f *Format) GetDefaultUnit(data float64) unit.Unit {
	if f.Type == Time {
		return unit.Second
	} else if f.Type == Bytes {
		return unit.Kb
	} else if f.Type == Percent {
		if data > 1 {
			return unit.Percent
		} else {
			return unit.PercentDecimal
		}
	} else {
		return unit.Default
	}
}

type Type uint8

const (
	Currency Type = iota
	Percent
	Number
	Time
	Bytes
)

func (d Type) IsPercent() bool {
	return d == Percent
}

func (d Type) Flag(u unit.Unit) string {
	return u.String()
}

func decimal(value float64, dec int) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%."+strconv.Itoa(dec)+"f", value), 64)
	return value
}

func decimalString(value float64, dec int) string {
	return fmt.Sprintf("%."+strconv.Itoa(dec)+"f", value)
}
