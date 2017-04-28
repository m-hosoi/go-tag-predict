package webtools

import (
	"context"
	"fmt"
	"go-tag-predict/cache"
	"go-tag-predict/fileutil"
	"os"
	"strings"
	"time"
)

func ExampleGet() {
	res, err := Get(context.Background(), "https://golang.org/", Options{})
	fmt.Println(err == nil)
	fmt.Println(res != nil)
	fmt.Println(res.Header.Get("Content-Type"))
	// Output:
	// true
	// true
	// text/html; charset=utf-8
}
func ExampleGetTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	res, err := Get(ctx, "https://golang.org/", Options{})
	fmt.Println(err != nil)
	fmt.Println(res == nil)
	fmt.Println(strings.Contains(err.Error(), ""))

	// Output:
	// true
	// true
	// true
}
func ExampleGetWithCache() {
	// no cache
	res, err := GetWithCache(context.Background(), "https://golang.org/", Options{}, CacheOptions{})
	fmt.Println(err == nil)
	fmt.Println(res != nil)
	fmt.Println(res.Header.Get("Content-Type"))

	d := cache.GetDefaultCacheDir()
	fmt.Println(strings.HasPrefix(d, "/tmp/")) // 安全装置
	fmt.Println(fileutil.Exist(d))

	// キャッシュがない状態
	res, err = GetWithCache(context.Background(), "https://golang.org/", Options{}, CacheOptions{CacheExpire: 1 * time.Second})
	fmt.Println(err == nil)
	fmt.Println(res != nil)
	fmt.Println(res.Header.Get("Content-Type"))

	fmt.Println(fileutil.Exist(d))

	// キャッシュがある状態
	res, err = GetWithCache(context.Background(), "https://golang.org/", Options{}, CacheOptions{CacheExpire: 1 * time.Second})
	fmt.Println(err == nil)
	fmt.Println(res != nil)
	fmt.Println(res.Header.Get("Content-Type"))

	// Cleanup
	os.RemoveAll(d)
	fmt.Println(fileutil.Exist(d))

	// Output:
	// true
	// true
	// text/html; charset=utf-8
	// true
	// false
	// true
	// true
	// text/html; charset=utf-8
	// true
	// true
	// true
	// text/html; charset=utf-8
	// false
}
