package monitor

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/service"
)

func InitRouter(prefix string, srv service.List) *context.App {

	app := context.NewApp()
	route := app.Group(prefix)

	// show dashboard
	route.GET("/dashboard/:dashboard_name", auth.Middleware(db.GetConnection(srv)), ShowDashboard)

	// refresh gauge
	route.POST("/refresh/:gauge_id/gauge/dashboard/:dashboard_name", auth.Middleware(db.GetConnection(srv)), Refresh)

	return app
}
