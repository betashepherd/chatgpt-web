package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// FileExists 文件是否存在
func FileExists(f string) bool {
	if _, err := os.Stat(f); err == nil {
		return true
	}
	return false
}

// CreateFile 创建文件
func CreateFile(path string) (*os.File, error) {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}
	return os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
}

// MakeFileDir 创建文件目录
func MakeFileDir(filename string) error {
	dir := filepath.Dir(filename)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

// CalcFileSize 获取文件大小
func CalcFileSize(filename string) int64 {
	fi, err := os.Stat(filename)
	if err != nil {
		logrus.WithError(err).Error("stat file failed")
		return -1
	}
	return fi.Size()
}

// FileCopy 文件复制
func FileCopy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return 0, err
	}
	defer destination.Close()

	nBytes, err := io.Copy(destination, source)

	return nBytes, err
}

// ParseFilename 解析出文件路径、文件名、扩展名(带.)
func ParseFilename(filename string) (string, string, string) {
	p := filepath.Dir(filename)
	b := filepath.Base(filename)
	ext := filepath.Ext(b)
	name := strings.TrimSuffix(b, ext)
	return p, name, ext
}

// DesenseName 姓名脱敏
func DesenseName(name string) string {
	nameRune := []rune(name)
	l := len(nameRune)
	if l <= 1 {
		return name
	} else if l == 2 {
		return string(nameRune[:1]) + "*"
	} else {
		s := string(nameRune[:1])
		e := string(nameRune[l-1:])
		return s + strings.Repeat("*", l-2) + e
	}
}

// DesenseNumber 证件号脱敏
func DesenseNumber(number string) string {
	numberRune := []rune(number)
	l := len(numberRune)
	if l <= 1 {
		return number
	} else if l == 2 {
		return string(numberRune[:1]) + "*"
	} else if l <= 5 {
		return string(numberRune[:1]) + strings.Repeat("*", l-2) + string(numberRune[l-1:])
	} else if l <= 7 {
		return string(numberRune[:2]) + strings.Repeat("*", l-4) + string(numberRune[l-2:])
	} else if l <= 9 {
		return string(numberRune[:3]) + strings.Repeat("*", l-6) + string(numberRune[l-3:])
	} else {
		return string(numberRune[:4]) + strings.Repeat("*", l-8) + string(numberRune[l-4:])
	}
}

// MaxLenRune 取字符串的指定最大长度的子串，按照字符计算，不区分中英文
func MaxLenRune(str string, maxLen int) string {
	strRune := []rune(str)
	if len(strRune) <= maxLen {
		return str
	}
	return string(strRune[:maxLen])
}

// MaxLenString 取字符串的指定最大长度的子串，按照字节计算，中文有可能被分割
func MaxLenString(str string, maxLen int) string {
	if len(str) <= maxLen {
		return str
	}
	return string(str[:maxLen])
}
