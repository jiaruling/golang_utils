package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func InitDir(path []string, delFileSuffix []string) {
	// 初始化文件操作对象
	f := NewFile()
	// 创建文件夹
	f.CreateDir(path)
	// 判断文件夹中的文件类型, 有相同类型的文件则删除
	for _, item := range path {
		for _, suffix := range delFileSuffix {
			if logList, err := f.WalkDir(item, suffix); err != nil {
				fmt.Println("扫描指定文件夹失败")
				os.Exit(1)
			} else {
				// 删除指定文件
				if len(logList) > 0 {
					f.DeletedFile(logList)
				}
			}
		}
	}
}

type File struct{}

func NewFile() *File {
	return &File{}
}

// 删除文件
func (f *File) DeletedFile(fileList []string) {
	for _, item := range fileList {
		if f.Exists(item) {
			_ = os.Remove(item)
		}
	}
}

// 判断所给路径文件/文件夹是否存在
func (f *File) Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if !os.IsExist(err) {
			return false
		}
	}
	return true
}

// 判断所给路径是否为文件夹
func (f *File) IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func (f *File) IsFile(path string) bool {
	return !f.IsDir(path)
}

// 获取文件夹下指定文件
func (f *File) WalkDir(dir, suffix string) (files []string, err error) {
	files = []string{}
	err = filepath.Walk(dir, func(fname string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			//忽略目录
			return nil
		}

		if len(suffix) == 0 || strings.HasSuffix(strings.ToLower(fi.Name()), suffix) {
			//文件后缀匹配
			files = append(files, fname)
		}

		return nil
	})

	return files, err
}

// 创建文件夹
func (f *File) CreateDir(pathList []string) {
	for _, path := range pathList {
		if !f.Exists(path) {
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				fmt.Println("创建目录"+path+"失败: ", err.Error())
				os.Exit(1)
			}
		}
	}
}
