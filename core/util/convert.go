package util

import "time"

// ParseTime 字符串转时间，默认格式 2006-01-02 15:04:05
func ParseTime(s string, layout ...string) (time.Time, error) {
	local, _ := time.LoadLocation("Local")
	if len(layout) > 0 {
		l := layout[0]
		t, err := time.ParseInLocation(l, s, local)
		return t, err
	}
	t, err := time.ParseInLocation("2006-01-02 15:04:05", s, local)
	return t, err
}

// TimeFormat 将时间格式化为字符串
func TimeFormat(t time.Time, layout ...string) string {
	s := ""
	if t.IsZero() {
		return ""
	}
	if len(layout) == 0 {
		s = t.In(time.Local).Format("2006-01-02 15:04:05")
	} else {
		s = t.In(time.Local).Format(layout[0])
	}
	return s
}
