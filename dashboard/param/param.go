package param

type Param map[string]interface{}

func (p Param) Get(key string) interface{} {
	if value, ok := p[key]; ok {
		return value
	}
	return ""
}

func (p Param) GetString(key string) string {
	if value, ok := p[key]; ok {
		return value.(string)
	}
	return ""
}

func (p Param) GetStringArray(key string) []string {
	if value, ok := p[key]; ok {
		return value.([]string)
	}
	return []string{}
}

func (p Param) GetInt64(key string) int64 {
	if value, ok := p[key]; ok {
		return value.(int64)
	}
	return 0
}

func (p Param) Combine(param Param) Param {
	for key, value := range param {
		exist := false
		for key2 := range p {
			if key == key2 {
				exist = true
				break
			}
		}
		if !exist {
			p[key] = value
		}
	}
	return p
}
