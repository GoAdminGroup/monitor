package monitor

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
)

func InitRouter(prefix string) *context.App {

	app := context.NewApp()
	route := app.Group(prefix)
	route.GET("/monitor/index", auth.Middleware, IndexHandler)

	return app
}
