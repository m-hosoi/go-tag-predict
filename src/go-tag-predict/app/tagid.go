package app

import (
	"bufio"
	"go-tag-predict/lambda"
	"io"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

// TagID : stringのタグをintに変換する
type TagID interface {
	GetID(tag string) int
	GetIDs(ar []string) []int
	Serialize(w io.Writer)
	GetReverse() map[int]string
}
type tagID struct {
	mutex  sync.Mutex
	seed   int
	tagMap map[string]int
}

// NewTagID : コンストラクタ
func NewTagID() TagID {
	return &tagID{tagMap: make(map[string]int, 1024*10)}
}

// LoadTagID : TagIDをロード
func LoadTagID(r io.Reader) (TagID, error) {
	scanner := bufio.NewScanner(r)
	tagMap := make(map[string]int, 1024*10)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		ar := strings.SplitN(scanner.Text(), "\t", 2)
		i, err := strconv.Atoi(ar[0])
		if err != nil {
			return nil, errors.Wrap(err, "parse error. line: "+strconv.Itoa(lineNo))
		}
		tagMap[ar[1]] = i
	}
	if err := scanner.Err(); err != nil {
		return nil, errors.WithStack(err)
	}
	return &tagID{tagMap: tagMap}, nil
}

func (t *tagID) getID(tag string) int {
	tag = strings.Replace(tag, "\n", "", -1)
	if i, ok := t.tagMap[tag]; ok {
		return i
	}
	t.seed++
	t.tagMap[tag] = t.seed
	return t.seed
}

// GetID : stringのタグをintに変換する
func (t *tagID) GetID(tag string) int {
	defer t.mutex.Unlock()
	t.mutex.Lock()
	return t.getID(tag)
}

// GetIDs : stringのタグをintに変換する
func (t *tagID) GetIDs(ar []string) []int {
	defer t.mutex.Unlock()
	t.mutex.Lock()

	return lambda.MapStringInt(ar, func(s string) int {
		return t.getID(s)
	})
}

// Serialize : シリアライズ
func (t *tagID) Serialize(w io.Writer) {
	bw := bufio.NewWriterSize(w, 1024*100)
	defer bw.Flush()
	for k, v := range t.tagMap {
		bw.WriteString(strconv.Itoa(v))
		bw.WriteString("\t")
		bw.WriteString(k)
		bw.WriteString("\n")
	}
}

// GetReverse : int => string変換表を取得する
func (t *tagID) GetReverse() map[int]string {
	idMap := make(map[int]string, len(t.tagMap))
	for k, v := range t.tagMap {
		idMap[v] = k
	}
	return idMap
}
