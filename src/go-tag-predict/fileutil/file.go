package fileutil

import (
	"go-tag-predict/lambda"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Exist : 存在チェック
func Exist(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// GetGopath : GOPATHを取得する
func GetGopath() []string {
	s := os.Getenv("GOPATH")
	if s == "" {
		panic("GOPATH is empty")
	}
	return lambda.MapString(strings.Split(s, ":"), strings.TrimSpace)
}

// FindFilePath : カレントディレクトリかGOPATHから指定のファイルを探す
func FindFilePath(name string) string {
	s, err := filepath.Abs(findFilePath(name))
	if err != nil {
		return name
	}
	return s
}
func findFilePath(name string) string {
	if Exist(name) {
		return name
	}
	for _, dir := range GetGopath() {
		s := path.Join(dir, name)
		if Exist(s) {
			return s
		}
	}
	return name
}
