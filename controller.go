package monitor

import (
	"bytes"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/go-admin/modules/logger"
	"github.com/GoAdminGroup/go-admin/modules/menu"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/constant"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/monitor/dashboard"
	"github.com/GoAdminGroup/monitor/dashboard/param"
	"net/http"
	"strconv"
)

func (m *Monitor) ShowDashboard(ctx *context.Context) {

	var (
		user  = auth.Auth(ctx)
		chart types.Panel
	)

	dashboardName := ctx.Query("dashboard_name")

	board := dashboard.Get(dashboardName)

	content, err := board.GetContent([]param.Param{
		{"interval": ctx.Query("interval")},
	})

	if err != nil {
		logger.Error("SetPageContent", err)
		alert := template.Get(config.GetTheme()).
			Alert().
			SetTitle(template.HTML(`<i class="icon fa fa-warning"></i> ` + language.Get("error") + `!`)).
			SetTheme("warning").SetContent(template.HTML(err.Error())).GetContent()
		chart = types.Panel{
			Content:     alert,
			Description: language.GetFromHtml("error"),
			Title:       language.GetFromHtml("error"),
		}
	} else {
		chart = types.Panel{
			Content:     content,
			Description: board.GetDescription(),
			Title:       board.GetTitle(),
		}
	}

	tmpl, tmplName := template.Get(config.GetTheme()).GetTemplate(ctx.Headers(constant.PjaxHeader) == "true")

	ctx.AddHeader("Content-Type", "text/html; charset=utf-8")

	buf := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(buf, tmplName, types.NewPage(types.NewPageParam{
		User:   user,
		Menu:   menu.GetGlobalMenu(user, m.Conn).SetActiveClass(config.URLRemovePrefix(ctx.Path())),
		Panel:  chart.GetContent(config.IsProductionEnvironment()),
		Assets: template.GetComponentAssetListsHTML(),
	}))
	if err != nil {
		logger.Error("ShowDashboard", err)
	}
	ctx.WriteString(buf.String())
}

func (m *Monitor) Refresh(ctx *context.Context) {
	dashboardName := ctx.Query("dashboard_name")

	board := dashboard.Get(dashboardName)

	id, _ := strconv.Atoi(ctx.Query("chart_id"))

	chart := board.GetChart(id)

	var p = param.NewFromFormValue(ctx.PostForm())

	data, err := chart.GetData(dashboard.NewWithChartType(dashboard.Graph).Combine(p))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "ok",
		"data": data,
	})
}
