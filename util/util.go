package util

import (
	"fmt"
	"github.com/google/shlex"
	"io"
	"os"
	"strings"
)

func ParseCmd(s string) []string {
	args, err := shlex.Split(s)
	if err != nil {
		panic(err)
	}
	return args
}

// 判断文件是否存在
func PathExist(path string) (bool, error) {
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
func CreateDir(path string) bool {
	exist, err := PathExist(path)
	if err != nil {
		panic(err)
	}
	if exist {
		fmt.Println("project has existed")
		return false
	} else {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			panic(err)
		}
		return true
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
	ok, _ := PathExist(dst)
	for ok {
		dst = fmt.Sprintf("%s_%s", dst, "cp")
		ok, _ = PathExist(dst)
	}
	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

func CopyExes(srcpath, dstpath string) {
	files, err := os.ReadDir(srcpath)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.IsDir() {
			newsrcpath := fmt.Sprintf("%s\\%s", srcpath, file.Name())
			CopyExes(newsrcpath, dstpath)
		} else {
			if strings.HasSuffix(file.Name(), ".exe") {
				newsrcpath := fmt.Sprintf("%s\\%s", srcpath, file.Name())
				newdstpath := fmt.Sprintf("%s\\%s", dstpath, file.Name())
				CopyFile(newsrcpath, newdstpath)
			}
		}
	}
}
