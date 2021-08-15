package mysql

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Datetime 可以适应Null值的时间类型
type Datetime struct {
	sql.NullTime
}

// NewDatetime 创建日期时间类型
func NewDatetime(t time.Time) Datetime {
	valid := true
	if t.IsZero() {
		valid = false
	}
	return Datetime{NullTime: struct {
		Time  time.Time
		Valid bool
	}{Time: t, Valid: valid}}
}

// MarshalJSON 转为json
func (v *Datetime) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Time)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON 转为Datetime
func (v Datetime) UnmarshalJSON(data []byte) error {
	var t time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	if t != (time.Time{}) {
		v.Valid = true
		v.Time = t
	} else {
		v.Valid = false
	}
	return nil
}
