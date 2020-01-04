package prometheus

import (
	"context"
	"errors"
	"github.com/GoAdminGroup/monitor/dashboard"
	"github.com/GoAdminGroup/monitor/dashboard/param"
	apicenter "github.com/prometheus/client_golang/api"
	prometheusapi "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"time"
)

type Prometheus struct {
	api      prometheusapi.API
	timezone *time.Location
}

const (
	InitKeyAddr     = "addr"
	InitKeyTimeZone = "timezone"
	DSKey           = "prometheus"

	QueryStatementKey = "query_statement"
	QueryStartTimeKey = "query_time_start"
	QueryEndTimeKey   = "query_time_end"
	QueryTimeStepKey  = "query_time_step"
)

func init() {
	dashboard.RegisterDS(DSKey, new(Prometheus))
}

func InitParam(addr string, timezone *time.Location) param.Param {
	return param.Param{InitKeyAddr: addr, InitKeyTimeZone: timezone}
}

type QueryConfig struct {
	Statements []string
	StartTime  int64
	EndTime    int64
	TimeStep   int64
}

func NewQueryConfig(statement ...string) QueryConfig {
	return QueryConfig{
		Statements: statement,
	}
}

func (q QueryConfig) AddStatement(statement string) QueryConfig {
	q.Statements = append(q.Statements, statement)
	return q
}

func (q QueryConfig) AddStartTime(start int64) QueryConfig {
	q.StartTime = start
	return q
}

func (q QueryConfig) AddEndTime(end int64) QueryConfig {
	q.EndTime = end
	return q
}

func (q QueryConfig) AddTimeStep(step int64) QueryConfig {
	q.TimeStep = step
	return q
}

func QueryParam(cfg QueryConfig) param.Param {
	return param.Param{
		QueryStatementKey: cfg.Statements,
		QueryStartTimeKey: cfg.StartTime,
		QueryEndTimeKey:   cfg.EndTime,
		QueryTimeStepKey:  cfg.TimeStep,
	}
}

func (p *Prometheus) GetName() string {
	return DSKey
}

func (p *Prometheus) Init(param param.Param) error {
	client, err := apicenter.NewClient(apicenter.Config{
		Address: param.GetString(InitKeyAddr),
	})
	if err != nil {
		return err
	}
	p.api = prometheusapi.NewAPI(client)
	p.timezone = param.Get(InitKeyTimeZone).(*time.Location)
	return nil
}

func (p *Prometheus) GetData(param param.Param) (dashboard.ChartData, error) {
	ctx := context.TODO()

	chartType := dashboard.GetChartType(param)

	var (
		values = make([]model.Value, 0)
		querys = param.GetStringArray(QueryStatementKey)
	)

	if chartType.IsGraph() {
		startTime := getTime(param.GetInt64(QueryStartTimeKey))
		endTime := getTime(param.GetInt64(QueryEndTimeKey))
		stepTime := time.Duration(param.GetInt64(QueryTimeStepKey)) * time.Second
		for _, query := range querys {
			// TODO: handle error
			val, _ := p.api.QueryRange(ctx, query, prometheusapi.Range{
				Start: startTime,
				End:   endTime,
				Step:  stepTime,
			})
			values = append(values, val)
		}
		data := dashboard.NewGraphData()

		durType := getTimeDurationType(startTime, endTime)

		for index, value := range values {
			v, _ := value.(model.Matrix)
			data.AddYAxis()
			for _, i := range v {
				for _, j := range i.Values {
					if index == 0 {
						data.AddLabels(p.getTimeLabel(j.Timestamp.Time(), durType))
					}
					data.AddYAxisPoint(index, float64(j.Value))
				}
			}
		}

		return data, nil
	}

	val, _ := p.api.Query(ctx, querys[0], time.Now())
	values = append(values, val)

	value, _ := values[0].(model.Vector)

	if len(value) == 0 {
		return dashboard.SingleStatData(0), errors.New("no data")
	}

	return dashboard.SingleStatData(value[0].Value), nil
}

func getTime(v int64) time.Time {
	return time.Unix(v, 0)
}

type DurType uint8

const (
	DurTypeMinute DurType = iota
	DurTypeHour
	DurTypeDay
	DurTypeMonth
	DurTypeYear
)

func getTimeDurationType(start, end time.Time) DurType {
	yearStart := start.Year()
	yearEnd := end.Year()

	if yearStart != yearEnd {
		return DurTypeYear
	}

	monethStart := start.Month()
	monethEnd := end.Month()

	if monethStart != monethEnd {
		return DurTypeMonth
	}

	dayStart := start.Day()
	dayEnd := end.Day()

	if dayStart != dayEnd {
		return DurTypeDay
	}

	hourStart := start.Hour()
	hourEnd := end.Hour()

	if hourStart != hourEnd {
		return DurTypeHour
	}

	return DurTypeMinute
}

func (p *Prometheus) getTimeLabel(t time.Time, durType DurType) string {
	switch durType {
	case DurTypeMinute:
		return t.In(p.timezone).Format("04:05")
	case DurTypeHour:
		return t.In(p.timezone).Format("15:04")
	case DurTypeDay:
		return t.In(p.timezone).Format(`01/02
15:04`)
	case DurTypeMonth:
		return t.In(p.timezone).Format(`01/02
15:04`)
	case DurTypeYear:
		return t.In(p.timezone).Format("2006/01/02")
	}
	panic("wrong duration type")
}
