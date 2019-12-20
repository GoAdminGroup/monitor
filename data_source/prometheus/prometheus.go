package prometheus

import (
	"context"
	"github.com/GoAdminGroup/monitor/dashboard"
	"github.com/GoAdminGroup/monitor/dashboard/param"
	apicenter "github.com/prometheus/client_golang/api"
	prometheusapi "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"strings"
	"time"
)

type Prometheus struct {
	api prometheusapi.API
}

const (
	InitKey = "addr"
	DSKey   = "prometheus"

	QueryStatementKey = "query_statement"
	QueryTypeKey      = "query_type"
	QueryTimeKey      = "query_time"
	QueryEndTimeKey   = "query_time_end"
	QueryTimeStepKey  = "query_time_step"
)

func init() {
	dashboard.RegisterDS(DSKey, new(Prometheus))
}

func InitParam(addr string) param.Param {
	return param.Param{InitKey: addr}
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
		QueryTimeKey:      cfg.StartTime,
		QueryEndTimeKey:   cfg.EndTime,
		QueryTimeStepKey:  cfg.TimeStep,
	}
}

func (p *Prometheus) GetName() string {
	return DSKey
}

func (p *Prometheus) Init(param param.Param) error {
	client, err := apicenter.NewClient(apicenter.Config{
		Address: param.GetString(InitKey),
	})
	if err != nil {
		return err
	}
	p.api = prometheusapi.NewAPI(client)
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
		startTime := getTime(param.GetInt64(QueryTimeKey))
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
	} else {
		val, _ := p.api.Query(ctx, querys[0], time.Now())
		values = append(values, val)
	}

	var cstZone = time.FixedZone("UTC", 8*3600)

	if chartType.IsGraph() {
		data := dashboard.NewGraphData()

		for index, value := range values {
			v, _ := value.(model.Matrix)
			data.AddYAxis()
			for _, i := range v {
				for _, j := range i.Values {
					if index == 0 {
						data.AddLabels(strings.Replace(j.Timestamp.Time().In(cstZone).Format("2006-01-02 15:04:05"), " ", `
`, -1))
					}
					data.AddYAxisPoint(index, float64(j.Value))
				}
			}

		}

		return data, nil
	} else {
		value, _ := values[0].(model.Vector)
		return dashboard.SingleStatData(value[0].Value), nil
	}
}

func getTime(v int64) time.Time {
	return time.Unix(v, 0)
}
