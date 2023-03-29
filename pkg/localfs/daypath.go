package localfs

import (
	"chatgpt-web/pkg/utils"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type DayPath struct {
	rootPath string
	urlHost  string
}

func NewDayPath(root, host string) *DayPath {
	p := DayPath{
		rootPath: strings.TrimRight(root, "/"),
		urlHost:  strings.TrimRight(host, "/"),
	}
	err := os.MkdirAll(p.rootPath, 0755)
	if err != nil {
		logrus.WithError(err).Error("init local fs failed")
		return nil
	}

	p.ClearTempPath()

	return &p
}

func (p *DayPath) ClearBefore(days int) {
	ls, err := ioutil.ReadDir(p.rootPath)
	if err != nil {
		return
	}

	t := utils.Now().AddDate(0, 0, -days)
	for _, l := range ls {
		if !l.IsDir() {
			continue
		}
		if len(l.Name()) != 8 {
			continue
		}
		tm, err := time.Parse("20060102", l.Name())
		if err != nil {
			continue
		}

		if tm.Unix() <= t.Unix() {
			_ = os.RemoveAll(utils.CombinePath(p.rootPath, l.Name()))
		}
	}
}

func (p *DayPath) ClearTempPath() {
	tmp, _ := p.GetTmpPath(false)
	_ = os.RemoveAll(tmp)
}

func (p *DayPath) RemovePath(pp string) {
	_ = os.RemoveAll(pp)
}

func (p *DayPath) GetPath(autoCreate bool, sub ...string) (string, error) {
	dayRoot := utils.CombinePath(p.rootPath, utils.TimeFormatDay(time.Now()))
	pp := utils.CombinePath(dayRoot, utils.CombinePath(sub...))
	if autoCreate {
		err := os.MkdirAll(pp, 0755)
		if err != nil {
			return "", err
		}
	}
	return pp, nil
}

func (p *DayPath) GetTmpPath(autoCreate bool, sub ...string) (string, error) {
	pp := utils.CombinePath(p.rootPath, "tmp", utils.CombinePath(sub...))
	if autoCreate {
		err := os.MkdirAll(pp, 0755)
		if err != nil {
			return "", err
		}
	}
	return pp, nil
}

// GetRelativePath 获取相对路径，即去掉 rootPath 的路径
func (p *DayPath) GetRelativePath(path string) string {
	return strings.TrimPrefix(path, p.rootPath+"/")
}

// GetLocalPath 获取本地路径，即加上 rootPath 的路径
func (p *DayPath) GetLocalPath(path string) string {
	if strings.HasPrefix(path, p.rootPath) {
		return path
	}
	return utils.CombinePath(p.rootPath, path)
}

// GetFullURL 获取本地文件对应的 URL
func (p *DayPath) GetFullURL(path string) string {
	if path == "" {
		return ""
	}
	return utils.CombinePath(p.urlHost, path)
}

// GetFullURLWithTime 获取本地文件对应的 URL
func (p *DayPath) GetFullURLWithTime(path string, tm time.Time) string {
	if path == "" {
		return ""
	}
	return utils.CombinePath(p.urlHost, path) + fmt.Sprintf("?t=%d", tm.Unix())
}

func (p *DayPath) SaveInterface(name string, obj interface{}, sub ...string) (string, error) {
	now := time.Now()
	dat, err := json.Marshal(obj)
	logrus.WithField("day-path", 1).Info("marshal obj ", sub, " cost ", time.Now().Sub(now).String())
	if err != nil {
		return "", err
	}
	return p.SaveDataFile(name, dat, sub...)
}

func (p *DayPath) SaveDataFile(name string, data []byte, sub ...string) (string, error) {
	pp, err := p.GetPath(true, sub...)
	if err != nil {
		return "", err
	}

	if name == "" {
		name = uuid.NewV4().String()
	}
	filename := utils.CombinePath(pp, name)

	now := time.Now()
	err = ioutil.WriteFile(filename, data, 0644)
	logrus.WithField("day-path", 1).Info("save file ", filename, " size ", FormatSize(int64(len(data))), " cost ", time.Now().Sub(now).String())
	if err != nil {
		return "", err
	}

	return filename, nil
}

func (p *DayPath) ReadDataFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func (p *DayPath) ReadInterface(filename string, obj interface{}) error {
	now := time.Now()
	dat, err := ioutil.ReadFile(filename)
	logrus.WithField("day-path", 1).Info("read file ", filename, " size ", FormatSize(int64(len(dat))), " cost ", time.Now().Sub(now).String())
	if err != nil {
		return err
	}
	now = time.Now()
	err = json.Unmarshal(dat, obj)
	logrus.WithField("day-path", 1).Info("unmarshal obj ", filename, " size ", FormatSize(int64(len(dat))), " cost ", time.Now().Sub(now).String())
	return err
}

func FormatSize(fileSize int64) string {
	if fileSize < 1024 {
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { // if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}
