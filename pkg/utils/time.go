package utils

import "time"

var defLocal, _ = time.LoadLocation("Asia/Shanghai")

// GetDefaultLocal 默认时区，中国
func GetDefaultLocal() *time.Location {
	return defLocal
}

// GetTime 根据各字段获取时间对象，时区默认中国
func GetTime(year, month, day, hour, min, sec int) time.Time {
	return time.Date(year, time.Month(month), day, hour, min, sec, 0, defLocal)
}

// GetMorningTime 获取 t 在当天中0点的时间对象，时区默认中国
func GetMorningTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, defLocal)
}

// GetMorningUnixTime 获取 t 在当天中0点的时间戳，时区默认中国
func GetMorningUnixTime(t time.Time) int64 {
	return GetMorningTime(t).Unix()
}

// GetSecondsOfDay 获取 t 在当天中的秒数，时区默认中国
func GetSecondsOfDay(t time.Time) int64 {
	return t.Unix() - GetMorningUnixTime(t)
}

// TimeFormat 自定义时间格式化
func TimeFormat(t time.Time, layout string) string {
	return t.In(defLocal).Format(layout)
}

// TimeFormatStandard format: 2006-01-02 15:04:05
func TimeFormatStandard(t time.Time) string {
	return t.In(defLocal).Format("2006-01-02 15:04:05")
}

// TimeFormatStandardMilli format: 2006-01-02 15:04:05.000
func TimeFormatStandardMilli(t time.Time) string {
	return t.In(defLocal).Format("2006-01-02 15:04:05.999")
}

// TimeFormatStamp format: 20060102150405
func TimeFormatStamp(t time.Time) string {
	return t.In(defLocal).Format("20060102150405")
}

// TimeFormatDay format: 20060102
func TimeFormatDay(t time.Time) string {
	return t.In(defLocal).Format("20060102")
}

// TimeFormatStampMilli format: 20060102150405.000
func TimeFormatStampMilli(t time.Time) string {
	return t.In(defLocal).Format("20060102150405.000")
}

// TimeFormatRFC3339 format: 2006-01-02T15:04:05Z07:00
func TimeFormatRFC3339(t time.Time) string {
	return t.In(defLocal).Format("2006-01-02T15:04:05Z07:00")
}

// TimeFormatRFC3339Milli format: 2006-01-02T15:04:05.000Z07:00
func TimeFormatRFC3339Milli(t time.Time) string {
	return t.In(defLocal).Format("2006-01-02T15:04:05.000Z07:00")
}

// GetTimeByUnixTime 通过时间戳获取时间对象，时区默认中国
func GetTimeByUnixTime(t int64) time.Time {
	return time.Unix(t, 0).In(defLocal)
}

// GetTimeByUnixTimeMillisecond 通过毫秒时间戳获取时间对象，时区默认中国
func GetTimeByUnixTimeMillisecond(t int64) time.Time {
	return time.Unix(t/1000, (t%1000)*1000000).In(defLocal)
}

// ParseTime 按照自定义格式解析时间字符串
func ParseTime(ts, layout string) (time.Time, error) {
	return time.ParseInLocation(layout, ts, defLocal)
}

// ParseTimeRFC3339 格式：2006-01-02T15:04:05.000Z07:00
func ParseTimeRFC3339(s string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05.000Z07:00", s)
}

// ParseTimeStandard 格式：2006-01-02 15:04:05
func ParseTimeStandard(s string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", s)
}

// Now 获取当前时间，时区默认中国
func Now() time.Time {
	return time.Now().In(defLocal)
}
