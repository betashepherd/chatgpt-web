package util

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"reflect"
	"regexp"
	"runtime"
	"strconv"

	"os"
	"path/filepath"
	"strings"
	"time"
)

var defLocal, _ = time.LoadLocation("Asia/Shanghai")

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

func GetDayString() string {
	return time.Now().In(defLocal).Format("20060102")
}

func GetOffsetTime(t time.Time, s int) time.Time {
	str := fmt.Sprintf("%ds", s)
	d, _ := time.ParseDuration(str)
	return t.Add(d)
}

func GetTimeStrByTime(t time.Time) string {
	return t.In(defLocal).Format("2006-01-02 15:04:05")
}

func GetDateByTime(t time.Time) string {
	str := t.In(defLocal).Format("20060102")
	return str
}

func GetTimeByUnixTime(t int64) time.Time {
	return time.Unix(t, 0)
}

func GetTimeByStrUnixTime(t string) (time.Time, error) {
	ts, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return GetTimeByUnixTime(ts), nil
}

func GetTimeByStrUnixMilliTime(t string) (time.Time, error) {
	ts, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return GetTimeByUnixTime(ts / 1000), nil
}

func GetTime(year, month, day, hour, min, sec int) time.Time {
	return time.Date(year, time.Month(month), day, hour, min, sec, 0, defLocal)
}

func GetMorningUnixTime(t time.Time) int64 {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, defLocal).Unix()
}

func GetNightUnixTime(t time.Time) int64 {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, defLocal).Unix()
}

func ParseTimeStrToTime(s string) time.Time {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", s, defLocal)
	if err != nil {
		return time.Time{}
	}
	return t
}

// TimeFormatStampMilli format: 20060102150405.000
func TimeFormatStampMilli(t time.Time) string {
	return t.In(defLocal).Format("20060102150405000")
}

func ParseTimeToStrWithLayout(layout string, t time.Time) string {
	return t.In(defLocal).Format(layout)
}

func ParseTimeStrToTimeWithLayout(layout, s string) (time.Time, error) {
	t, err := time.ParseInLocation(layout, s, defLocal)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func GetUnixTimeByTime(t time.Time) int64 {
	return t.Unix()
}

func GetCurrentTime() time.Time {
	return time.Now().In(defLocal)
}

func FileExists(f string) bool {
	if _, err := os.Stat(f); err == nil {
		return true
	}
	return false
}

func CreateFile(path string) (*os.File, error) {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return nil, err
		}
	}
	return os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
}

func GetMd5String(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}

func Sha256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func IsIntInSlice(arr []int, a int) bool {
	if len(arr) == 0 {
		return false
	}
	for _, ra := range arr {
		if ra == a {
			return true
		}
	}
	return false
}

func IsStrInSlice(arr []string, s string) bool {
	if len(arr) == 0 {
		return false
	}
	for _, ra := range arr {
		if ra == s {
			return true
		}
	}
	return false
}

func Base64Decode(data *string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(*data)
}

func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// GetRequestSETime 获取请求的时间范围，精度到天
// st, et 精度到秒的时间戳
func GetRequestSETime(st, et int64) (time.Time, time.Time, error) {
	if st == 0 {
		st = GetCurrentTime().Unix()
	}
	if et == 0 {
		et = GetCurrentTime().Unix()
	}
	if st > et {
		return time.Now(), time.Now(), errors.New("结束时间需大于开始时间")
	}

	stime := GetTimeByUnixTime(st)
	etime := GetTimeByUnixTime(et)

	loc, _ := time.LoadLocation("Asia/Shanghai")
	s := time.Date(stime.Year(), stime.Month(), stime.Day(), 0, 0, 0, 0, loc)
	e := time.Date(etime.Year(), etime.Month(), etime.Day(), 23, 59, 59, 0, loc)
	return s, e, nil
}

// 求并集
func Union(slice1, slice2 []string) []string {
	m := make(map[string]int)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 0 {
			slice1 = append(slice1, v)
		}
	}
	return slice1
}

// 求交集
func Intersect(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}

// 去重
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

func CalcPasswordResetLeftSeconds(lastUpdate time.Time) int64 {
	lastDay := lastUpdate.AddDate(0, 3, 0)
	return lastDay.Unix() - GetCurrentTime().Unix()
}

func CombinePath(paths ...string) (p string) {
	for _, s := range paths {
		if p == "" {
			p = strings.TrimRight(s, "/")
		} else {
			p += "/" + strings.Trim(s, "/")
		}
	}
	return p
}

func GetFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

func Strings2Ints(ss []string) []int {
	is := []int{}

	for _, v := range ss {
		vi, _ := strconv.Atoi(v)
		is = append(is, vi)
	}
	return is
}

func Ints2Strings(is []int) []string {
	ss := []string{}
	for _, i := range is {
		ss = append(ss, fmt.Sprintf("%d", i))
	}

	return ss
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

func Mt_rand(min, max int64) int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int63n(max-min+1) + min
}

func StringsUnion(slice1, slice2 []string) []string {
	m := map[string]int{}
	for _, v1 := range slice1 {
		m[v1] = 1
	}

	for _, v2 := range slice2 {
		if _, ok := m[v2]; !ok {
			m[v2] = 1
		}
	}

	union := []string{}
	for s, _ := range m {
		union = append(union, s)
	}

	return union
}

func StringsIntersect(slice1, slice2 []string) []string {
	m := map[string]int{}
	it := []string{}
	for _, v1 := range slice1 {
		m[v1] = 1
	}

	for _, v2 := range slice2 {
		if _, ok := m[v2]; ok {
			it = append(it, v2)
		}
	}
	return it
}

func FileGetContent(fileName string) ([]string, error) {
	data, err := ioutil.ReadFile("data/update_events")
	if err != nil {
		return nil, err
	}
	contents := strings.Split(string(data), "\n")
	ss := []string{}
	for _, line := range contents {
		if strings.TrimSpace(line) == "" {
			continue
		}
		ss = append(ss, line)
	}

	return ss, nil
}

//06:01 ->
func HourMin2Sec(t string) int64 {
	t = strings.Replace(t, "-", "", -1)
	t = strings.Replace(t, " ", "", -1)
	t = strings.Replace(t, ":", "", -1)
	t = strings.Replace(t, "/", "", -1)
	it, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		return 0
	}

	h := it / 100
	m := it % 100
	if h < 0 || h > 23 {
		return 0
	}
	if m < 0 || m > 59 {
		return 0
	}
	return h*3600 + m*60
}

func TimeDiffDay(st, et time.Time) int {
	st = time.Date(st.Year(), st.Month(), st.Day(), 0, 0, 0, 0, defLocal)
	et = time.Date(et.Year(), et.Month(), et.Day(), 0, 0, 0, 0, defLocal)
	return int(et.Sub(st).Hours() / 24)
}

func DayRange(st, et time.Time) []string {
	diff := TimeDiffDay(st, et)
	ds := []string{st.Format("20060102")}
	if diff > 1 {
		for i := 1; i <= diff; i++ {
			ds = append(ds, st.AddDate(0, 0, i).Format("20060102"))
		}
	}
	return ds
}

func CopyFile(srcName, dstName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return 0, err
	}
	defer src.Close()

	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return 0, err
	}
	defer dst.Close()

	return io.Copy(dst, src)
}

func Gzip(in string) []byte {
	var out bytes.Buffer
	gz := gzip.NewWriter(&out)
	gz.Write([]byte(in))
	gz.Flush()
	gz.Close()
	return out.Bytes()
}

func UnGzip(bin []byte) string {
	var in bytes.Buffer
	var out bytes.Buffer
	in.Write(bin)
	r, _ := gzip.NewReader(&in)
	defer r.Close()
	io.Copy(&out, r)
	return out.String()
}

// email verify
func VerifyEmailFormat(email string) bool {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// mobile verify
func VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}
