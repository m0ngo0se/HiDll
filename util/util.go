package util

import (
	"fmt"
	"github.com/google/shlex"
	"io"
	"os"
)

func ParseCmd(s string) []string {
	args, err := shlex.Split(s)
	if err != nil {
		panic(err)
	}
	return args
}

// 判断文件夹是否存在
func HasDir(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 创建文件夹
func CreateDir(path string) {
	exist, err := HasDir(path)
	if err != nil {
		panic(err)
	}
	if exist {
		fmt.Println("project has existed")
	} else {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}

func CopyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}
