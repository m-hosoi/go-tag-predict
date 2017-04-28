package pinboard

import (
	"bytes"
	"context"
	"go-tag-predict/webtools"
	"io"
	"io/ioutil"
	"strings"

	"golang.org/x/net/html/charset"

	xpp "github.com/mmcdole/goxpp"
	"github.com/pkg/errors"
)

// Post : ブックマークデータ
type Post struct {
	Title string
	Href  string
	Tags  []string
}

// PostIterator : ブックマークデータを逐次取得する
type PostIterator interface {
	// Next : ブックマークデータを逐次取得する
	// データが無くなったら *Post == nil && error == nilを返す
	Next() (*Post, error)
}

type postIterator struct {
	p *xpp.XMLPullParser
}

// GetAll : ユーザーのブックマークを全て取得する
//
// 結果データが膨大になるのを考慮して、[]Postではなく、逐次処理するイテレータを返す
//
// Example:
//   ctx := context.Background()
//   itr, err := pinboard.GetAll(ctx, pinboardAPIToken, webtools.CacheOptions{CacheExpire: 24 * time.Hour, CacheDir: cacheDir})
//   for {
//     post, err := itr.Next()
//     if err != nil {
//       return err
//     }
//     if post == nil {
//       break
//     }
//     fmt.Println(post.Title)
//   }
func GetAll(ctx context.Context, apiToken string, co webtools.CacheOptions) (PostIterator, error) {
	res, err := webtools.GetWithCache(
		ctx, "https://api.pinboard.in/v1/posts/all?auth_token="+apiToken,
		webtools.Options{},
		co)
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
	if bytes.Equal(body, []byte("API requires authentication")) {
		return nil, errors.New(string(body))
	}
	if len(body) == 0 {
		return nil, errors.New("bad request")
	}
	return parseAllData(bytes.NewReader(body))
}

// LoadFile : ファイルからブックマークを取得する
func LoadFile(filePath string) (PostIterator, error) {
	body, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return parseAllData(bytes.NewReader(body))
}

func parseAllData(r io.Reader) (PostIterator, error) {
	p := xpp.NewXMLPullParser(r, true, newReaderLabel)

	err := find(p, xpp.StartTag)
	if err != nil {
		return nil, err
	}
	if "posts" != p.Name {
		return nil, errors.New("bad data")
	}

	return &postIterator{p: p}, nil
}

func (itr *postIterator) Next() (*Post, error) {
	p := itr.p
	for {
		err := find(p, xpp.StartTag)
		if err != nil {
			break
		}
		if "post" != p.Name {
			break
		}
		href := p.Attribute("href")
		description := p.Attribute("description")
		shared := p.Attribute("shared")
		tag := p.Attribute("tag")
		if shared == "no" { // 非公開のブックマークを除外
			continue
		}
		if href == "" || description == "" {
			continue
		}
		p := Post{Title: description, Href: href, Tags: strings.Split(tag, " ")}
		return &p, nil
	}
	return nil, nil
}

// github.com/mmcdole/gofeed/internal/shared/NewReaderLabel
func newReaderLabel(label string, input io.Reader) (io.Reader, error) {
	conv, err := charset.NewReaderLabel(label, input)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return conv, nil
}
func find(p *xpp.XMLPullParser, target xpp.XMLEventType) error {
	for {
		e, err := p.Next()
		if err != nil {
			return errors.WithStack(err)
		}
		if e == target {
			return nil
		}
		if e == xpp.EndDocument {
			return io.EOF
		}
	}
}
