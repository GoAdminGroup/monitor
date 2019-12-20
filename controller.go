package monitor

import (
	"bytes"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/go-admin/modules/logger"
	"github.com/GoAdminGroup/go-admin/modules/menu"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/constant"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/monitor/dashboard"
	"github.com/GoAdminGroup/monitor/dashboard/param"
)

func ShowDashboard(ctx *context.Context) {

	var (
		user  = auth.Auth(ctx)
		chart types.Panel
	)

	dashboardName := ctx.Query("dashboard_name")

	board := dashboard.Get(dashboardName)

	content, err := board.GetContent([]param.Param{})

	if err != nil {
		logger.Error("SetPageContent", err)
		alert := template.Get(config.Theme).
			Alert().
			SetTitle(template.HTML(`<i class="icon fa fa-warning"></i> ` + language.Get("error") + `!`)).
			SetTheme("warning").SetContent(template.HTML(err.Error())).GetContent()
		chart = types.Panel{
			Content:     alert,
			Description: language.Get("error"),
			Title:       language.Get("error"),
		}
	} else {
		chart = types.Panel{
			Content:     content,
			Description: board.GetDescription(),
			Title:       board.GetTitle(),
		}
	}

	tmpl, tmplName := template.Get(config.Theme).GetTemplate(ctx.Headers(constant.PjaxHeader) == "true")

	ctx.AddHeader("Content-Type", "text/html; charset=utf-8")

	buf := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(buf, tmplName, types.NewPage(user,
		*(menu.GetGlobalMenu(user, connection).SetActiveClass(config.URLRemovePrefix(ctx.Path()))),
		chart, config, template.GetComponentAssetListsHTML()))
	if err != nil {
		logger.Error("ShowDashboard", err)
	}
	ctx.WriteString(buf.String())
}

func Refresh(ctx *context.Context) {

}
