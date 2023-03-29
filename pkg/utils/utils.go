package utils

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"reflect"
	"runtime"
)

// CalcStringMD5 计算字符串md5
func CalcStringMD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}

// CalcFileMD5 计算文件md5
func CalcFileMD5(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		logrus.WithError(err).Error("open file failed")
		return ""
	}
	defer f.Close()

	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		logrus.WithError(err).Error("copy file failed")
		return ""
	}

	has := md5hash.Sum(nil)
	return fmt.Sprintf("%x", has)
}

// Base64Decode base64 编码
func Base64Decode(image *string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(*image)
}

// Base64Encode base64 解码
func Base64Encode(imageByte []byte) string {
	return base64.StdEncoding.EncodeToString(imageByte)
}

// GetFuncName 获取当前代码所属函数名称
func GetFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

// InArray will search element inside array with any type.
// Will return boolean and index for matched element.
// True and index more than 0 if element is exist.
// needle is element to search, haystack is slice of value to be search.
func InArray(needle interface{}, haystack interface{}) (bool, int) {
	exists := false
	index := -1

	switch reflect.TypeOf(haystack).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(haystack)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(needle, s.Index(i).Interface()) == true {
				index = i
				exists = true
				break
			}
		}
	}

	return exists, index
}

func IsInArray(elem interface{}, array interface{}) bool {
	exists := false

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(elem, s.Index(i).Interface()) == true {
				exists = true
				break
			}
		}
	}

	return exists
}

func MinInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func MinUint64(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}

func MinInt64(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func MaxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func MaxInt64(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func MaxUint64(x, y uint64) uint64 {
	if x > y {
		return x
	}
	return y
}
