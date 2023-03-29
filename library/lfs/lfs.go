package lfs

import "chatgpt-web/pkg/localfs"

var DataFs *localfs.DayPath

func Init(root, host string) {
	DataFs = localfs.NewDayPath(root, host)
	if DataFs == nil {
		panic("cannot create data path")
	}
}
