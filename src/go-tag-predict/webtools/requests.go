package webtools

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"go-tag-predict/cache"
	"net/http"
	"net/http/httputil"
	"path"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Options : リクエストオプション
type Options struct {
	UserAgent string
	Referer   string
}

// CacheOptions : Cacheオプション
type CacheOptions struct {
	CacheExpire time.Duration
	CacheDir    string
}

// Get : Getリクエスト
func Get(ctx context.Context, rawurl string, o Options) (*http.Response, error) {
	req, err := http.NewRequest("GET", rawurl, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if o.UserAgent != "" {
		req.Header.Set("User-Agent", o.UserAgent)
	}
	if o.Referer != "" {
		req.Header.Set("Referer", o.Referer)
	}
	return Request(ctx, req, o)
}

// GetWithCache : キャッシュ付きGetリクエスト
func GetWithCache(ctx context.Context, rawurl string, o Options, co CacheOptions) (*http.Response, error) {
	if co.CacheExpire == 0 {
		return Get(ctx, rawurl, o)
	}
	hash := sha256.New()
	hash.Write([]byte(rawurl + o.UserAgent))

	// 単一フォルダ内のファイルが増えすぎないように、階層化する
	cacheName := "webtools.request/" + strings.Join(strings.SplitN(hex.EncodeToString(hash.Sum(nil)), "", 4), "/")
	d := co.CacheDir
	if d == "" {
		d = cache.GetDefaultCacheDir()
	}
	var response *http.Response
	ar, err := cache.WithFileCache(func() ([]byte, error) {
		res, err := Get(ctx, rawurl, o)
		if err != nil {
			return nil, err
		}
		response = res
		if ar, err := serializeResponse(res); err != nil {
			return nil, errors.WithStack(err)
		} else {
			return ar, nil
		}
	}, co.CacheExpire, path.Join(d, cacheName))
	if err != nil {
		if response != nil {
			if response.Body != nil {
				response.Body.Close()
			}
		}
		return nil, err
	}
	if response != nil {
		return response, nil
	}
	return deserializeResponse(ar)
}

// Request : HTTPリクエストを送信する
func Request(ctx context.Context, req *http.Request, o Options) (*http.Response, error) {
	tr := &http.Transport{
		DisableKeepAlives: true, // keep alive onだとgoroutineがleakする
	}
	client := &http.Client{
		Transport: tr,
	}

	res, err := client.Do(req.WithContext(ctx))
	if err != nil {
		if res != nil { // エラーでもResponseが返ることがある
			if res.Body != nil {
				res.Body.Close()
			}
		}
		return nil, errors.WithStack(err)
	}
	if res.StatusCode >= 400 {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, errors.New(res.Status)
	}
	return res, nil
}

func serializeResponse(res *http.Response) ([]byte, error) {
	return httputil.DumpResponse(res, true)
}

func deserializeResponse(ar []byte) (*http.Response, error) {
	r := bufio.NewReader(bytes.NewReader(ar))
	return http.ReadResponse(r, nil)
}
