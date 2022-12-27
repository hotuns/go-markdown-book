//go:build darwin

package utils

import (
	"errors"
	"os"
	"runtime"
	"syscall"
)

type stat_t struct {
	CreateAt int64
	UpdateAt int64
}

func GetSysTime(path string) (stat_t, error) {
	finfo, err := os.Stat(path)
	if err != nil {
		return stat_t{}, err
	}

	var info stat_t

	// 是否是macos
	if runtime.GOOS == "darwin" {
		fattr, ok := finfo.Sys().(*syscall.Stat_t)
		if !ok {
			return stat_t{}, errors.New("not ok")
		}
		info = stat_t{
			CreateAt: fattr.Ctimespec.Sec,
			UpdateAt: fattr.Mtimespec.Sec,
		}
	} else {
		info = stat_t{
			CreateAt: finfo.ModTime().Unix(),
			UpdateAt: finfo.ModTime().Unix(),
		}
	}

	return info, nil
}
