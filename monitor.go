package monitor

import (
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/service"
	"github.com/GoAdminGroup/go-admin/plugins"
	"github.com/GoAdminGroup/monitor/dashboard"
)

type Monitor struct {
	*plugins.Base
}

const Name = "monitor"

func NewMonitor() *Monitor {
	return &Monitor{
		Base: &plugins.Base{PlugName: Name},
	}
}

func (m *Monitor) AddDashboard(name string, gen dashboard.Gen) *Monitor {
	dashboard.Add(name, gen())
	return m
}

func (m *Monitor) InitPlugin(srv service.List) {
	// DO NOT DELETE
	m.InitBase(srv)

	m.Conn = db.GetConnection(srv)

	m.App = m.initRouter(srv)
	addToLanguagePkg()
}
