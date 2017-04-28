package app

import (
	"bytes"
	"context"
	"go-tag-predict/webtools"
	"io/ioutil"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	readability "github.com/mauidude/go-readability"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
)

// LoadWebContent : WebページからHTMLタグを除去したテキストを取得する
func LoadWebContent(ctx context.Context, rawurl string, cacheDir string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	// TODO: 相手のサーバーへ負荷を掛けすぎないように、同一ドメインへのリクエストは間隔を空ける
	res, err := webtools.GetWithCache(
		ctx, rawurl,
		webtools.Options{},
		webtools.CacheOptions{CacheExpire: 1 * time.Hour, CacheDir: cacheDir})
	// res, err := webtools.Get(
	// ctx, rawurl,
	// webtools.Options{})
	if err != nil {
		if res != nil {
			panic(err)
		}
		return "", err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	if -1 == strings.Index(res.Header.Get("Content-Type"), "html") {
		return "", errors.New(res.Header.Get("Content-Type") + " not supported")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.WithStack(err)
	}
	// TODO: utf-8以外のエンコードは、考慮していない
	enc := webtools.DetectTextEncode(&res.Header, body)
	body, err = webtools.DecodeText(body, enc)
	if err != nil {
		return "", errors.WithStack(err)
	}

	ry, err := readability.NewDocument(string(body))
	if err != nil {
		return "", errors.WithStack(err)
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(ry.Content()))
	text := doc.Find("body").Text()
	return strings.TrimSpace(text), nil
}

// LoadFeed : RSS/Atomフィードを取得する
func LoadFeed(ctx context.Context, rawurl string, cacheDir string) (*gofeed.Feed, error) {
	res, err := webtools.GetWithCache(
		ctx, rawurl,
		webtools.Options{},
		webtools.CacheOptions{CacheExpire: 1 * time.Hour, CacheDir: cacheDir})
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	fp := gofeed.NewParser()
	feed, err := fp.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return feed, nil
}
