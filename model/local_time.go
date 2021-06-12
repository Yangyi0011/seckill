package model

import (
	"database/sql/driver"
	"time"
)

const (
	TimeFormat = "2006-01-02 15:04:05"
)

// LocalTime 自定义 LocalTime 转换规则，以支持 yyyy-MM-dd HH:mm:ss 格式
type LocalTime time.Time

func (t *LocalTime) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 2 {
		*t = LocalTime(time.Time{})
		return
	}

	now, err := time.Parse(`"`+TimeFormat+`"`, string(data))
	*t = LocalTime(now)
	return
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TimeFormat)
	b = append(b, '"')
	return b, nil
}

func (t LocalTime) Value() (driver.Value, error) {
	if t.String() == "0001-01-01 00:00:00" {
		return nil, nil
	}
	return []byte(time.Time(t).Format(TimeFormat)), nil
}

func (t *LocalTime) Scan(v interface{}) error {
	tTime, _ := time.Parse("2006-01-02 15:04:05 +0800 CST", v.(time.Time).String())
	*t = LocalTime(tTime)
	return nil
}

func (t LocalTime) String() string {
	return time.Time(t).Format(TimeFormat)
}

func (t LocalTime) ZeroValue() LocalTime {
	zero, _ := time.Parse(TimeFormat, "0001-01-01 00:00:00")
	return LocalTime(zero)
}

// Unix 转为 time.Time.Unix()
// 如果 LocalTime 是直接从数据库查出来的，此时 LocalTime 相对于 time.Time 是 +8个小时的，
// 不能直接使用 time.Time(t).Unix()，需要做一下处理
func (t LocalTime) Unix() int64 {
	return time.Time(t).Add(-8*time.Hour).Unix()
}
