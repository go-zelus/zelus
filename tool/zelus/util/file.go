package util

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
)

// FileExisted 检查文件是否存在
func FileExisted(filename string) bool {
	existed := true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		existed = false
	}
	return existed
}

func CopyFile(src, dest string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	fi, _ := srcFile.Stat()
	perm := fi.Mode()

	destFile, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return 0, err
	}
	defer destFile.Close()

	return io.Copy(destFile, srcFile)
}

func copyFile(src, dst string, replaces []string) error {
	var err error
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	buf, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	var old string
	for i, next := range replaces {
		if i%2 == 0 {
			old = next
			continue
		}
		buf = bytes.ReplaceAll(buf, []byte(old), []byte(next))
	}
	return ioutil.WriteFile(dst, buf, srcInfo.Mode())
}
