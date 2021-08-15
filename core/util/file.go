package util

import "os"

// IsDir 判断路径是否是文件夹
func IsDir(p string) bool {
	s, err := os.Stat(p)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile 判断是否是文件
func IsFile(p string) bool {
	return !IsDir(p)
}

// Exists 判断路径/文件夹是否存在
func Exists(p string) bool {
	_, err := os.Stat(p)
	if err != nil && os.IsExist(err) {
		return true
	}
	return false
}
