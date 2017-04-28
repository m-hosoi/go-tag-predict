package cache

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"go-tag-predict/fileutil"
	"strings"
	"time"
)

func ExampleWithFileCache() {
	d := GetDefaultCacheDir()
	fmt.Println(strings.HasPrefix(d, "/tmp/")) // 安全装置
	fmt.Println(fileutil.Exist(d))

	p := path.Join(d, "filecache_test.cache")
	fmt.Println(fileutil.Exist(p))

	_, err := WithFileCache(func() ([]byte, error) {
		return nil, errors.New("test error")
	}, 1*time.Second, p)
	fmt.Println(err)

	res, err := WithFileCache(func() ([]byte, error) {
		return []byte("result 1"), nil
	}, 1*time.Second, p)
	fmt.Println(err)
	fmt.Println(string(res))
	fmt.Println(fileutil.Exist(p))

	// キャッシュ有効期限内
	res, err = WithFileCache(func() ([]byte, error) {
		return []byte("result 2"), nil
	}, 1*time.Second, p)
	fmt.Println(string(res))

	// キャッシュ有効期限切れ
	time.Sleep(1 * time.Second)
	res, err = WithFileCache(func() ([]byte, error) {
		return []byte("result 2"), nil
	}, 1*time.Second, p)
	fmt.Println(string(res))

	// Cleanup
	os.RemoveAll(d)
	fmt.Println(fileutil.Exist(d))
	// Output:
	// true
	// false
	// false
	// test error
	// <nil>
	// result 1
	// true
	// result 1
	// result 2
	// false
}
func ExampleWithFileCacheNull() {
	d := GetDefaultCacheDir()
	fmt.Println(strings.HasPrefix(d, "/tmp/")) // 安全装置
	fmt.Println(fileutil.Exist(d))

	p := path.Join(d, "filecache_test.cache")
	fmt.Println(fileutil.Exist(p))

	res, err := WithFileCache(func() ([]byte, error) {
		return nil, nil
	}, 1*time.Second, p)
	fmt.Println(err)
	fmt.Println(res == nil)

	fmt.Println(fileutil.Exist(p))

	// Cleanup
	os.RemoveAll(d)
	fmt.Println(fileutil.Exist(d))

	// Output:
	// true
	// false
	// false
	// <nil>
	// true
	// false
	// false
}
func ExampleRemoveAllExpired() {
	makeDummy := func(s string) {
		f, err := os.Create(s)
		defer f.Close()
		if err != nil {
			return
		}
		f.Write([]byte(s))
	}
	getFileCount := func(s string) int {
		ar, err := ioutil.ReadDir(s)
		if err != nil {
			return 0
		}
		return len(ar)
	}
	d := GetDefaultCacheDir()
	fmt.Println(strings.HasPrefix(d, "/tmp/")) // 安全装置
	fmt.Println(fileutil.Exist(d))
	_ = os.MkdirAll(d, 0700)
	fmt.Println(fileutil.Exist(d))

	makeDummy(path.Join(d, "tmp0"))
	makeDummy(path.Join(d, "tmp1"))
	time.Sleep(1 * time.Second)
	makeDummy(path.Join(d, "tmp2"))
	fmt.Println(getFileCount(d))

	RemoveAllExpired(d, 1*time.Second)
	fmt.Println(getFileCount(d))
	time.Sleep(1 * time.Second)
	RemoveAllExpired(d, 1*time.Second)
	fmt.Println(getFileCount(d))

	// Cleanup
	os.RemoveAll(d)
	fmt.Println(fileutil.Exist(d))
	// Output:
	// true
	// false
	// true
	// 3
	// 1
	// 0
	// false
}
