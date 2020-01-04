package dashboards

import (
	"github.com/GoAdminGroup/monitor/chart"
	"github.com/GoAdminGroup/monitor/chart/data"
	"github.com/GoAdminGroup/monitor/chart/data/unit"
	"github.com/GoAdminGroup/monitor/dashboard"
	"github.com/GoAdminGroup/monitor/data_source/prometheus"
	"time"
)

func NodeExporter() dashboard.Dashboard {
	board := dashboard.NewDefaultDashboard("Node Exporter", "prometheus node exporter")

	board.AddDataSource(prometheus.DSKey, prometheus.InitParam("http://localhost:9090", time.FixedZone("UTC", 8*3600)))

	h, _ := time.ParseDuration("-1h")
	start := time.Now().Add(h).Unix()
	end := time.Now().Unix()

	query7 := `sum(time() - node_boot_time_seconds{instance=~"localhost:9100"})`

	board.AddChartWithDSKey(prometheus.DSKey, chart.NewSingleStat(chart.SingleStatConfig{
		BaseConfig: chart.BaseConfig{
			Title:       "系统运行时间",
			Description: "123123",
			Size:        []int{2, 6},
			Height:      130,
			DataFormat: &data.Format{
				Type:     data.Time,
				Decimals: 1,
			},
			Param: prometheus.QueryParam(prometheus.NewQueryConfig(query7).
				AddStartTime(start).AddEndTime(end).
				AddTimeStep(15)),
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
			DataFormat: &data.Format{
				Type: data.Number,
			},
			Param: prometheus.QueryParam(prometheus.NewQueryConfig(query8).
				AddStartTime(start).AddEndTime(end).
				AddTimeStep(15)),
		},
		Color: "#b6b6c7",
	}))

	query9 := `sum(node_memory_MemTotal_bytes{instance=~"localhost:9100"})`

	board.AddChartWithDSKey(prometheus.DSKey, chart.NewSingleStat(chart.SingleStatConfig{
		BaseConfig: chart.BaseConfig{
			Title:       "内存总量",
			Description: "123123",
			Size:        []int{2, 6},
			Height:      130,
			DataFormat: &data.Format{
				Type:     data.Time,
				Decimals: 1,
			},
			Param: prometheus.QueryParam(prometheus.NewQueryConfig(query9).
				AddStartTime(start).AddEndTime(end).
				AddTimeStep(15)),
		},
		Color: "#b6b6c7",
	}))

	query10 := `avg(irate(node_cpu_seconds_total{instance=~"localhost:9100",mode="iowait"}[30m])) * 100`

	board.AddChartWithDSKey(prometheus.DSKey, chart.NewSingleStat(chart.SingleStatConfig{
		BaseConfig: chart.BaseConfig{
			Title:       "CPU iowait",
			Description: "123123",
			Size:        []int{2, 6},
			Height:      130,
			DataFormat: &data.Format{
				Type:     data.Time,
				Decimals: 1,
			},
			Param: prometheus.QueryParam(prometheus.NewQueryConfig(query10).
				AddStartTime(start).AddEndTime(end).
				AddTimeStep(15)),
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
			DataFormat: &data.Format{
				Type:     data.Percent,
				Decimals: 2,
			},
			Param: prometheus.QueryParam(prometheus.NewQueryConfig(query1, query2, query3).
				AddStartTime(start).AddEndTime(end).
				AddTimeStep(15)),
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
			DataFormat: &data.Format{
				Type:     data.Number,
				Decimals: 2,
			},
			Param: prometheus.QueryParam(prometheus.NewQueryConfig(query4, query5, query6).
				AddStartTime(start).AddEndTime(end).
				AddTimeStep(15)),
		},
		YLabels: []string{"1分钟", "5分钟", "15分钟"},
	}))

	query11 := `irate(node_network_receive_bytes_total{instance=~'localhost:9100',device!~'tap.*|veth.*|br.*|docker.*|virbr*|lo*'}[30m])*8`
	query12 := `irate(node_network_transmit_bytes_total{instance=~'localhost:9100',device!~'tap.*|veth.*|br.*|docker.*|virbr*|lo*'}[30m])*8`

	board.AddChartWithDSKey(prometheus.DSKey, chart.NewGraph(chart.GraphConfig{
		BaseConfig: chart.BaseConfig{
			Title:       "网络流量",
			Description: "123123",
			Size:        []int{6, 12},
			Height:      330,
			DataFormat: &data.Format{
				Type:     data.Bytes,
				DSUnit:   unit.Kb,
				Unit:     unit.Mb,
				Decimals: 1,
			},
			Param: prometheus.QueryParam(prometheus.NewQueryConfig(query11, query12).
				AddStartTime(time.Now().Add(2 * h).Unix()).AddEndTime(end).
				AddTimeStep(15)),
		},
		YLabels: []string{"上传", "下载"},
	}))

	query13 := `irate(node_disk_read_time_seconds_total{instance=~"localhost:9100"}[30m]) / irate(node_disk_reads_completed_total{instance=~"localhost:9100"}[30m])`
	//query14 := `irate(node_disk_write_time_seconds_total{instance=~"localhost:9100"}[30m]) / irate(node_disk_writes_completed_total{instance=~"localhost:9100"}[30m])`
	//query15 := `irate(node_disk_io_time_seconds_total{instance=~"localhost:9100"}[30m])`
	//query16 := `irate(node_disk_io_time_weighted_seconds_total{instance=~"localhost:9100"}[30m])`

	board.AddChartWithDSKey(prometheus.DSKey, chart.NewGraph(chart.GraphConfig{
		BaseConfig: chart.BaseConfig{
			Title:       "每次IO读写的耗时",
			Description: "123123",
			Size:        []int{6, 12},
			Height:      330,
			DataFormat: &data.Format{
				Type:     data.Time,
				DSUnit:   unit.Millisecond,
				Unit:     unit.Millisecond,
				Decimals: 1,
			},
			Param: prometheus.QueryParam(prometheus.NewQueryConfig(query13).
				AddStartTime(time.Now().Add(2 * h).Unix()).AddEndTime(end).
				AddTimeStep(15)),
		},
		YLabels: []string{"disk0_读取"},
	}))

	return board
}
