package dashboards

import (
	"github.com/GoAdminGroup/monitor/chart"
	"github.com/GoAdminGroup/monitor/dashboard"
	"github.com/GoAdminGroup/monitor/data_source/prometheus"
	"time"
)

func NodeExporter() dashboard.Dashboard {
	board := dashboard.NewDefaultDashboard("Node Exporter", "prometheus node exporter")

	board.AddDataSource(prometheus.DSKey, prometheus.InitParam("http://localhost:9090"))

	h, _ := time.ParseDuration("-1h")

	query7 := `sum(time() - node_boot_time_seconds{instance=~"localhost:9100"})`

	board.AddChartWithDSKey(prometheus.DSKey, chart.NewSingleStat(chart.SingleStatConfig{
		BaseConfig: chart.BaseConfig{
			Title:       "系统运行时间",
			Description: "123123",
			Size:        []int{2, 6},
			Height:      130,
			DataFormat: chart.Format{
				Type:     chart.Second,
				Decimals: 1,
			},
			Param: prometheus.QueryParam(prometheus.NewQueryConfig(query7). // , query2, query3, query4
				AddStartTime(time.Now().Add(h).Unix()).AddEndTime(time.Now().Unix()).
				AddTimeStep(60 * 5)),
		},
		Color: "#b6b6c7",
	}))

	query8 := `sum(count(node_cpu_seconds_total{instance=~"localhost:9100", mode='system'}) by (cpu))`

	board.AddChartWithDSKey(prometheus.DSKey, chart.NewSingleStat(chart.SingleStatConfig{
		BaseConfig: chart.BaseConfig{
			Title:       "CPU核数",
			Description: "123123",
			Size:        []int{2, 6},
			Height:      130,
			DataFormat: chart.Format{
				Type:     chart.Number,
			},
			Param: prometheus.QueryParam(prometheus.NewQueryConfig(query8). // , query2, query3, query4
				AddStartTime(time.Now().Add(h).Unix()).AddEndTime(time.Now().Unix()).
				AddTimeStep(60 * 5)),
		},
		Color: "#b6b6c7",
	}))

	board.AddChartWithDSKey(prometheus.DSKey, chart.NewSingleStat(chart.SingleStatConfig{
		BaseConfig: chart.BaseConfig{
			Title:       "内存总量",
			Description: "123123",
			Size:        []int{2, 6},
			Height:      130,
			DataFormat: chart.Format{
				Type:     chart.Second,
				Decimals: 1,
			},
			Param: prometheus.QueryParam(prometheus.NewQueryConfig(query7). // , query2, query3, query4
				AddStartTime(time.Now().Add(h).Unix()).AddEndTime(time.Now().Unix()).
				AddTimeStep(60 * 5)),
		},
		Color: "#b6b6c7",
	}))

	board.AddChartWithDSKey(prometheus.DSKey, chart.NewSingleStat(chart.SingleStatConfig{
		BaseConfig: chart.BaseConfig{
			Title:       "CPU iowait",
			Description: "123123",
			Size:        []int{2, 6},
			Height:      130,
			DataFormat: chart.Format{
				Type:     chart.Second,
				Decimals: 1,
			},
			Param: prometheus.QueryParam(prometheus.NewQueryConfig(query7). // , query2, query3, query4
				AddStartTime(time.Now().Add(h).Unix()).AddEndTime(time.Now().Unix()).
				AddTimeStep(60 * 5)),
		},
		Color: "#b6b6c7",
	}))

	query1 := `avg(irate(node_cpu_seconds_total{instance=~"localhost:9100",mode="system"}[30m])) by (instance)`
	query2 := `avg(irate(node_cpu_seconds_total{instance=~"localhost:9100",mode="user"}[30m])) by (instance)`
	query3 := `avg(irate(node_cpu_seconds_total{instance=~"localhost:9100",mode="iowait"}[30m])) by (instance)`
	//query4 := `1 - avg(irate(node_cpu_seconds_total{instance=~"localhost:9100",mode="idle"}[30m])) by (instance)`

	board.AddChartWithDSKey(prometheus.DSKey, chart.NewGraph(chart.GraphConfig{
		BaseConfig: chart.BaseConfig{
			Title:       "cpu使用率",
			Description: "123123",
			Size:        []int{4, 12},
			Height:      330,
			DataFormat: chart.Format{
				Type:     chart.Percent,
				Decimals: 2,
			},
			Param: prometheus.QueryParam(prometheus.NewQueryConfig(query1, query2, query3). // , query2, query3, query4
				AddStartTime(time.Now().Add(h).Unix()).AddEndTime(time.Now().Unix()).
				AddTimeStep(60 * 5)),
		},
		YLabels: []string{"system", "user", "iowait"},
	}))

	query4 := `node_load1{instance=~"localhost:9100"}`
	query5 := `node_load5{instance=~"localhost:9100"}`
	query6 := `node_load15{instance=~"localhost:9100"}`

	board.AddChartWithDSKey(prometheus.DSKey, chart.NewGraph(chart.GraphConfig{
		BaseConfig: chart.BaseConfig{
			Title:       "系统平均负载",
			Description: "123123",
			Size:        []int{4, 12},
			Height:      330,
			DataFormat: chart.Format{
				Type:     chart.Number,
				Decimals: 2,
			},
			Param: prometheus.QueryParam(prometheus.NewQueryConfig(query4, query5, query6). // , query2, query3, query4
				AddStartTime(time.Now().Add(h).Unix()).AddEndTime(time.Now().Unix()).
				AddTimeStep(60 * 5)),
		},
		YLabels: []string{"1分钟", "5分钟", "15分钟"},
	}))

	return board
}
