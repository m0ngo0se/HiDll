package util

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func ParseCmd(s string) []string {
	args, err := Split(s)
	if err != nil {
		println(err)
	}
	return args
}

// 判断文件是否存在
func PathExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// 创建文件夹
func CreateDir(path string) bool {
	exist := PathExist(path)
	if exist {
		fmt.Println("project has existed")
		return false
	} else {
		err := os.Mkdir(path, os.ModePerm)
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
	ok := PathExist(dst)
	for ok {
		dst = fmt.Sprintf("%s_%s", dst, "cp")
		ok = PathExist(dst)
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

func Writedata(path string, date string) {
	fd, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	fd.Write(append([]byte(date), '\n'))
	fd.Close()
}
