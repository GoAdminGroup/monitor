package main

import (
	c "github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/service"
	"github.com/GoAdminGroup/go-admin/plugins"
	e "github.com/GoAdminGroup/go-admin/plugins/example"
)

type Example struct {
	*plugins.Base
}

var Plugin = &Example{
	Base: &plugins.Base{PlugName: "example"},
}

var config c.Config

func (example *Example) InitPlugin(srv service.List) {
	config = c.Get()
	Plugin.App = e.InitRouter(config.Prefix(), srv)
	e.SetConfig(config)
}
