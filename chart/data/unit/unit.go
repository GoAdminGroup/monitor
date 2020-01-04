package unit

type Unit uint8

const (
	None Unit = iota

	Default

	Millisecond
	Second
	Minute
	Hour
	Day

	Kb
	Mb
	Gb

	PercentDecimal
	Percent
)

var (
	Time  = []Unit{Millisecond, Second, Minute, Hour, Day}
	Bytes = []Unit{Kb, Mb, Gb}
)

func (u Unit) String() string {
	switch u {
	case Default:
		return ""
	case Millisecond:
		return "ms"
	case Second:
		return "second"
	case Minute:
		return "minute"
	case Hour:
		return "hour"
	case Day:
		return "day"
	case Kb:
		return "kb"
	case Mb:
		return "mb"
	case Gb:
		return "gb"
	case Percent:
		return "%"
	case PercentDecimal:
		return ""
	}
	return ""
}

func (u Unit) Coefficient() float64 {
	switch u {
	case Default:
		return 1
	case Millisecond:
		return 0.001
	case Second:
		return 1
	case Minute:
		return 60
	case Hour:
		return 3600
	case Day:
		return 86400
	case Kb:
		return 1
	case Mb:
		return 1024
	case Gb:
		return 1048576
	case Percent:
		return 1
	case PercentDecimal:
		return 100
	}
	return 1
}

func (u Unit) IsTime() bool {
	return u == Second || u == Minute || u == Hour || u == Day
}

func (u Unit) IsBytes() bool {
	return u == Kb || u == Mb || u == Gb
}

func (u Unit) IsPercent() bool {
	return u == Percent || u == PercentDecimal
}

func (u Unit) IsValid() bool {
	return u.IsBytes() || u.IsTime() || u.IsDefault() || u.IsPercent()
}

func (u Unit) IsDefault() bool {
	return u == Default
}

func GetUnitFromData(data float64, unit Unit) Unit {
	if unit.IsTime() {
		return GetUnitFromUnits(data, Time)
	} else if unit.IsBytes() {
		return GetUnitFromUnits(data, Bytes)
	} else if unit.IsPercent() {
		return Percent
	} else {
		return Default
	}
}

func GetUnitFromUnits(data float64, units []Unit) Unit {
	for i, un := range units {
		if i == len(units)-1 {
			return un
		}
		if data <= un.Coefficient() {
			if i == 0 {
				return un
			}
			return units[i-1]
		}
	}
	return Default
}
