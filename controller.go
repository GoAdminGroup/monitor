package monitor

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/page"
	"github.com/GoAdminGroup/go-admin/template/types"
)

func IndexHandler(ctx *context.Context) {
	page.SetPageContent(ctx, auth.Auth(ctx), func(ctx interface{}) (types.Panel, error) {
		return types.Panel{
			Content:     "",
			Title:       "Server Scanner",
			Description: "statistics of server",
		}, nil
	})
}
