package monitor

import (
	"github.com/GoAdminGroup/go-admin/context"
	c "github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/plugins"
)

type Monitor struct {
	app *context.App
}

func NewMonitor() *Monitor {
	return Plug
}

var Plug = new(Monitor)

var config c.Config

func SetConfig(cfg c.Config) {
	config = cfg
}

func (scanner *Monitor) InitPlugin() {
	config = c.Get()
	Plug.app = InitRouter(config.Prefix())
}

func (scanner *Monitor) GetRequest() []context.Path {
	return scanner.app.Requests
}

func (scanner *Monitor) GetHandler(url, method string) context.Handlers {
	return plugins.GetHandler(url, method, scanner.app)
}
