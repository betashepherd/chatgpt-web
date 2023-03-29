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

type LocalPath struct {
	rootPath string
	urlHost  string
}

func NewLocalPath(root, host string) *LocalPath {
	p := LocalPath{
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

func (p *LocalPath) ClearTempPath() {
	tmp, _ := p.GetTmpPath(false)
	_ = os.RemoveAll(tmp)
}

func (p *LocalPath) RemovePath(pp string) {
	_ = os.RemoveAll(pp)
}

func (p *LocalPath) GetPath(autoCreate bool, sub ...string) (string, error) {
	pp := utils.CombinePath(p.rootPath, utils.CombinePath(sub...))
	if autoCreate {
		err := os.MkdirAll(pp, 0755)
		if err != nil {
			return "", err
		}
	}
	return pp, nil
}

func (p *LocalPath) GetTmpPath(autoCreate bool, sub ...string) (string, error) {
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
func (p *LocalPath) GetRelativePath(path string) string {
	return strings.TrimPrefix(path, p.rootPath+"/")
}

// GetLocalPath 获取本地路径，即加上 rootPath 的路径
func (p *LocalPath) GetLocalPath(path string) string {
	if strings.HasPrefix(path, p.rootPath) {
		return path
	}
	return utils.CombinePath(p.rootPath, path)
}

// GetFullURL 获取本地文件对应的 URL
func (p *LocalPath) GetFullURL(path string) string {
	if path == "" {
		return ""
	}
	return utils.CombinePath(p.urlHost, path)
}

// GetFullURLWithTime 获取本地文件对应的 URL
func (p *LocalPath) GetFullURLWithTime(path string, tm time.Time) string {
	if path == "" {
		return ""
	}
	return utils.CombinePath(p.urlHost, path) + fmt.Sprintf("?t=%d", tm.Unix())
}

func (p *LocalPath) SaveInterface(name string, obj interface{}, sub ...string) (string, error) {
	now := time.Now()
	dat, err := json.Marshal(obj)
	logrus.WithField("local-path", 1).Info("marshal obj ", sub, " cost ", time.Now().Sub(now).String())
	if err != nil {
		return "", err
	}
	return p.SaveDataFile(name, dat, sub...)
}

func (p *LocalPath) SaveDataFile(name string, data []byte, sub ...string) (string, error) {
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
	logrus.WithField("local-path", 1).Info("save file ", filename, " size ", FormatSize(int64(len(data))), " cost ", time.Now().Sub(now).String())
	if err != nil {
		return "", err
	}

	return filename, nil
}

func (p *LocalPath) ReadDataFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func (p *LocalPath) ReadInterface(filename string, obj interface{}) error {
	now := time.Now()
	dat, err := ioutil.ReadFile(filename)
	logrus.WithField("local-path", 1).Info("read file ", filename, " size ", FormatSize(int64(len(dat))), " cost ", time.Now().Sub(now).String())
	if err != nil {
		return err
	}
	now = time.Now()
	err = json.Unmarshal(dat, obj)
	logrus.WithField("local-path", 1).Info("unmarshal obj ", filename, " size ", FormatSize(int64(len(dat))), " cost ", time.Now().Sub(now).String())
	return err
}
