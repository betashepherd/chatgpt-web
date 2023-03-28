package types

import (
	"chatgpt-web/pkg/logger"
	"strconv"
)

func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

func UInt64ToString(num uint64) string {
	return strconv.FormatUint(num, 10)
}

func StringToInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		logger.Danger(err, "StringToInt Err")
	}

	return num
}
