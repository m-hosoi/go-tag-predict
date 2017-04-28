package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"go-tag-predict/fileutil"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

// WithFileCache : 簡易なファイルキャッシュ
// ※同名のキャッシュへの同時アクセスは想定しない
// Example:
//   res, err := cache.WithFileCache(func() ([]byte, error) {
//     return []byte("result 1"), nil
//   }, 1*time.Hour, path.Join(cache.GetDefaultCacheDir(), "cache.dat"))
func WithFileCache(proc func() ([]byte, error), expire time.Duration, cachePath string) ([]byte, error) {
	res, err := get(cachePath, expire)
	if err != nil {
		return nil, err
	}
	if res != nil {
		return res, nil
	}

	res, err = proc()
	if err != nil {
		return nil, err
	}

	err = put(cachePath, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func get(cachePath string, expire time.Duration) ([]byte, error) {
	fi, err := os.Stat(cachePath)
	if !os.IsNotExist(err) {
		if expire > time.Since(fi.ModTime()) {
			f, err := ioutil.ReadFile(cachePath)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			return f, nil
		}
		os.Remove(cachePath)
	}
	return nil, nil
}
func put(cachePath string, res []byte) error {
	if res == nil {
		if !fileutil.Exist(cachePath) {
			return nil
		}
		if err := os.Remove(cachePath); err != nil {
			return errors.WithStack(err)
		}
		return nil
	}
	dir := path.Dir(cachePath)
	if !fileutil.Exist(dir) {
		err := os.MkdirAll(dir, 0700)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	if err := ioutil.WriteFile(cachePath, res, 0600); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// RemoveAllExpired : 期限切れの全てのキャッシュファイルを削除する
func RemoveAllExpired(cacheDir string, expire time.Duration) error {
	if !fileutil.Exist(cacheDir) {
		return nil
	}
	ar, err := ioutil.ReadDir(cacheDir)
	if err != nil {
		return errors.WithStack(err)
	}
	for _, fi := range ar {
		if expire <= time.Since(fi.ModTime()) {
			os.Remove(path.Join(cacheDir, fi.Name()))
		}
	}
	return nil
}

var cacheDir string

// GetDefaultCacheDir : デフォルトのキャッシュ保存場所
func GetDefaultCacheDir() string {
	if cacheDir != "" {
		return cacheDir
	}
	seed, _ := filepath.Abs(os.Args[0])
	hash := sha256.New()
	hash.Write([]byte(seed))
	name := hex.EncodeToString(hash.Sum(nil))
	cacheDir = path.Join("/tmp/", "cache_"+name)
	return cacheDir
}
