package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// PrettyJson 格式化结构体
func PrettyJson(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		return string(b)
	}
	return ""
}

// RemoveDuplicate 字符串数组去重
func RemoveDuplicate(arr []string) []string {
	resArr := make([]string, 0)
	tmpMap := make(map[string]interface{})
	for _, val := range arr {
		// 判断主键为val的map是否存在
		if _, ok := tmpMap[val]; !ok {
			resArr = append(resArr, val)
			tmpMap[val] = nil
		}
	}

	return resArr
}

// CombinePath 地址/路径拼接，自动去掉多余的斜线，一般用于连接url或者unix/linux路径，不能去掉相对路径
func CombinePath(paths ...string) (p string) {
	for _, s := range paths {
		if p == "" {
			p = strings.TrimRight(s, "/")
		} else {
			s := strings.Trim(s, "/")
			if s != "" {
				p += "/" + s
			}
		}
	}
	return p
}

// MarshalToString 对 json.Marshal 的封装
func MarshalToString(v interface{}) string {
	val, _ := json.Marshal(v)
	return string(val)
}

// UnmarshalFromString 对 json.Unmarshal 的封装
func UnmarshalFromString(d string, v interface{}) {
	_ = json.Unmarshal([]byte(d), v)
}

// GetPeriodString 格式化时长为中文描述
func GetPeriodString(seconds int64) string {
	if seconds <= 60 {
		return fmt.Sprintf("%d秒", seconds)
	}
	min := float64(seconds) / 60
	if min <= 60 {
		return fmt.Sprintf("%.1f分", min)
	}
	hour := min / 60
	if hour <= 24 {
		return fmt.Sprintf("%.1f小时", hour)
	}
	days := hour / 24
	if days <= 30 {
		return fmt.Sprintf("%.1f天", days)
	}
	return fmt.Sprintf("%.0f天前", days)
}

// Camel2Case 驼峰字符串转小写字符串
func Camel2Case(name string) string {
	var buffer bytes.Buffer
	lastCap := false
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 && !lastCap {
				buffer.WriteByte('_')
			}
			buffer.WriteRune(unicode.ToLower(r))
			lastCap = true
		} else {
			buffer.WriteRune(r)
			lastCap = false
		}
	}
	return buffer.String()
}

func Str2Int(s string, def int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return i
}

func Str2Uint(s string, def uint) uint {
	i, err := strconv.ParseUint(s, 10, 0)
	if err != nil {
		return def
	}
	return uint(i)
}

func Str2Int64(s string, def int64) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return def
	}
	return i
}

func Str2Uint64(s string, def uint64) uint64 {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return def
	}
	return i
}

func Str2Float32(s string, def float32) float32 {
	i, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return def
	}
	return float32(i)
}

func Str2Float64(s string, def float64) float64 {
	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return def
	}
	return i
}
