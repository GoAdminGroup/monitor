package monitor

import (
	"github.com/GoAdminGroup/go-admin/context"
	c "github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/service"
	"github.com/GoAdminGroup/monitor/dashboard"
)

type Monitor struct {
	app *context.App
}

func NewMonitor() *Monitor {
	return Plug
}

var Plug = new(Monitor)

var (
	config     c.Config
	connection db.Connection
)

func SetConfig(cfg c.Config) {
	config = cfg
}

func (monitor *Monitor) AddDashboard(name string, gen dashboard.Gen) *Monitor {
	dashboard.Add(name, gen())
	return monitor
}

func (monitor *Monitor) InitPlugin(services service.List) {
	config = c.Get()
	Plug.app = InitRouter(config.Prefix(), services)
	connection = db.GetConnection(services)
	addToLanguagePkg()
}

func (monitor *Monitor) GetHandler() context.HandlerMap {
	return monitor.app.Handlers
}
