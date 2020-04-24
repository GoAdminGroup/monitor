package monitor

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/service"
)

func (m *Monitor) initRouter(srv service.List) *context.App {

	app := context.NewApp()
	route := app.Group(config.GetUrlPrefix())

	// show dashboard
	route.GET("/dashboard/:dashboard_name", auth.Middleware(db.GetConnection(srv)), m.ShowDashboard)

	// refresh gauge
	route.POST("/refresh/:chart_id/chart/dashboard/:dashboard_name", auth.Middleware(db.GetConnection(srv)), m.Refresh)

	return app
}
