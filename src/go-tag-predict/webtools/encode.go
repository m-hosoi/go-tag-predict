package webtools

import (
	"bytes"
	"go-tag-predict/lambda"
	"io/ioutil"
	"net/http"
	"strings"
	"unicode"

	"github.com/pkg/errors"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// TextEncode : 文字コード
type TextEncode int

const (
	// TextEncodeUnknown : 不明
	TextEncodeUnknown TextEncode = iota
	// TextEncodeUtf8 : UTF-8
	TextEncodeUtf8
	// TextEncodeShiftJIS : Shift_JIS
	TextEncodeShiftJIS
	// TextEncodeEUCJP : euc-jp
	TextEncodeEUCJP
	// TextEncodeISO2022JP : iso-2022-jp
	TextEncodeISO2022JP
)

// DetectTextEncode : HTTPヘッダーと本文から文字コードを取得
//
// ※高速化のために厳密なHTMLのパースは行わない
// 例）<!-- <meta charset="UTF-8" /> --> これも文字コードとして検出してしまう
func DetectTextEncode(header *http.Header, body []byte) TextEncode {
	if s := detectEncodeFromHeader(header); s != "" {
		return toTextEncode(s)
	}
	if header != nil && -1 == strings.Index(header.Get("Content-Type"), "html") {
		return TextEncodeUnknown
	}
	if s := detectEncodeFromBody(body); s != "" {
		return toTextEncode(s)
	}
	return TextEncodeUnknown
}
func detectEncodeFromBody(body []byte) string {
	if body == nil {
		return ""
	}
	body = bytes.ToLower(toHead(body))
	if s := findValueForKey(body, []byte("encoding")); s != nil {
		return string(s)
	}
	if s := findValueForKey(body, []byte("charset")); s != nil {
		return string(s)
	}
	return ""
}

func findValueForKey(s []byte, k []byte) []byte {
	i := bytes.Index(s, k)
	if i == -1 {
		return nil
	}
	s = bytes.TrimLeftFunc(s[i+len(k):], unicode.IsSpace)
	if len(s) == 0 {
		return nil
	}
	if s[0] != '=' {
		return nil
	}
	s = bytes.TrimLeftFunc(s[1:], unicode.IsSpace)
	if len(s) == 0 {
		return nil
	}
	if s[0] == '"' || s[0] == '\'' {
		s = s[1:]
	}
	j := bytes.IndexFunc(s, func(r rune) bool {
		if '0' <= r && r <= '9' {
			return false
		}
		if 'a' <= r && r <= 'z' {
			return false
		}
		if '_' == r || '-' == r {
			return false
		}
		return true
	})
	if j == -1 {
		return nil
	}
	return s[:j]
}

// チェック対象をbodyタグ前までに絞る
func toHead(s []byte) []byte {
	if i := bytes.Index(s, []byte("<b")); i != -1 {
		return s[:i]
	} else if i := bytes.Index(s, []byte("<B")); i != -1 {
		return s[:i]
	}
	return s
}
func detectEncodeFromHeader(header *http.Header) string {
	if header == nil {
		return ""
	}
	return mimeToEncode(header.Get("Content-Type"))
}
func mimeToEncode(s string) string {
	for _, line := range lambda.MapString(strings.Split(s, ";"), func(line string) string {
		return strings.ToLower(strings.TrimSpace(line))
	}) {
		if !strings.HasPrefix(line, "charset") {
			continue
		}
		kv := lambda.MapString(strings.SplitN(line, "=", 2), strings.TrimSpace)
		if len(kv) != 2 {
			continue
		}
		if kv[0] != "charset" {
			continue
		}
		return kv[1]
	}
	return ""
}
func toTextEncode(s string) TextEncode {
	if s == "utf-8" {
		return TextEncodeUtf8
	}
	if s == "utf8" {
		return TextEncodeUtf8
	}
	if s == "shift_jis" {
		return TextEncodeShiftJIS
	}
	if s == "sjis" {
		return TextEncodeShiftJIS
	}
	if s == "shift-jis" {
		return TextEncodeShiftJIS
	}
	if s == "cp932" {
		return TextEncodeShiftJIS
	}
	if s == "eucjp" {
		return TextEncodeEUCJP
	}
	if s == "euc-jp" {
		return TextEncodeEUCJP
	}
	if s == "iso-2022-jp" {
		return TextEncodeISO2022JP
	}
	return TextEncodeUnknown
}

// DecodeText : 文字列をデコードする
func DecodeText(body []byte, enc TextEncode) ([]byte, error) {
	var t transform.Transformer
	if enc == TextEncodeUtf8 || enc == TextEncodeUnknown {
		return body, nil
	} else if enc == TextEncodeShiftJIS {
		t = japanese.ShiftJIS.NewDecoder()
	} else if enc == TextEncodeEUCJP {
		t = japanese.EUCJP.NewDecoder()
	} else if enc == TextEncodeISO2022JP {
		t = japanese.ISO2022JP.NewDecoder()
	} else {
		return nil, errors.Errorf("invalid text encode: %v", enc)
	}
	r := transform.NewReader(bytes.NewReader(body), t)
	if body, err := ioutil.ReadAll(r); err != nil {
		return nil, errors.WithStack(err)
	} else {
		return body, nil
	}
}
