package monitor

import (
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/monitor/constant"
)

var langs = map[string]map[string]string{
	"cn": {
		join("second"): "秒",
		join("minute"): "分钟",
		join("hour"):   "小时",
		join("day"):    "天",
	},
	"en": {

	},
}

func join(key string) string {
	return constant.Scope + "." + key
}

func addToLanguagePkg() {
	for key, lang := range langs {
		for k, v := range lang {
			language.Lang[key][k] = v
		}
	}
}
