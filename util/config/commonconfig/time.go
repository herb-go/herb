package commonconfig

import "time"

var defaultDateFormat = "2006-01-02"
var defaultDatetimeFormat = "2006-01-02 15:04:05"
var defaultTimeFormat = "15:04:05"

type TimeConfig struct {
	Timezone       string
	TimeFormat     string
	DateFormat     string
	DatetimeFormat string
	location       *time.Location
}

func (c *TimeConfig) TimeInLocation(t time.Time) time.Time {
	if c.location == nil {
		if c.Timezone == "" {
			c.location = time.Local
		} else {
			var err error
			c.location, err = time.LoadLocation(c.Timezone)
			if err != nil {
				panic(err)
			}
		}
	}
	return t.In(c.location)
}
func (c *TimeConfig) Date(t time.Time) string {
	localTime := c.TimeInLocation(t)
	if c.DateFormat == "" {
		return localTime.Format(defaultDateFormat)
	}
	return localTime.Format(c.DateFormat)
}

func (c *TimeConfig) Time(t time.Time) string {
	localTime := c.TimeInLocation(t)

	if c.TimeFormat == "" {
		return localTime.Format(defaultTimeFormat)
	}
	return localTime.Format(c.TimeFormat)
}

func (c *TimeConfig) Datetime(t time.Time) string {
	localTime := c.TimeInLocation(t)

	if c.DatetimeFormat == "" {
		return localTime.Format(defaultDatetimeFormat)
	}
	return localTime.Format(c.DatetimeFormat)
}
