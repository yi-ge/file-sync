package utils

import (
	"os"
	"runtime"
)

// CurrentFile - 获取当前的文件路径
func CurrentFile() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	return file
}

// CurrentDir - 获取当前执行文件的目录
func CurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

// MakeDirIfNotExist - 如果文件夹不存在则创建
func MakeDirIfNotExist(path string) {
	if !IsExist(path) {
		os.MkdirAll(path, 0755)
	}
}

// IsExist -判断文件/文件夹是否存在  存在返回 true 不存在返回false
func IsExist(filename string) bool {
	var exist = true

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}

	return exist
}
