package param

import "strconv"

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
		if v, ok := value.(int64); ok {
			return v
		}
		if v, ok := value.(string); ok {
			vv, _ := strconv.ParseInt(v, 10, 64)
			return vv
		}
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

func NewFromFormValue(f map[string][]string) Param {
	var p = make(Param)
	for k, v := range f {
		p[k] = v[0]
	}
	return p
}
