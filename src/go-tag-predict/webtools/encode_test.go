package webtools

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func ExampleDetectTextEncode() {
	h := http.Header{}
	fmt.Println(DetectTextEncode(&h, nil))
	h.Set("content-type", "text/html; charset=utf-8")
	fmt.Println(DetectTextEncode(&h, nil))
	h.Set("content-type", "text/html; charset=utf8")
	fmt.Println(DetectTextEncode(&h, nil))

	fmt.Println(DetectTextEncode(nil, []byte(`
	<?xml version="1.0" encoding="utf-8"?>
	<HTML>
	<body>test</body>
	</html>
`)))
	fmt.Println(DetectTextEncode(nil, []byte(`
	<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=Shift_JIS" />
	</head>
	</html>
`)))
	fmt.Println(DetectTextEncode(nil, []byte(`
	<html>
	<head>
		<meta charset="UTF-8" />
	</head>
	</html>
`)))

	// Output:
	// 0
	// 1
	// 1
	// 1
	// 2
	// 1
}

func ExampleDecodeText() {
	url := "http://www.soumu.go.jp/"
	res, err := Get(context.Background(), url, Options{})
	fmt.Println(err)
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(err)
	if res.Body != nil {
		defer res.Body.Close()
	}
	enc := DetectTextEncode(&res.Header, body)
	fmt.Println(enc == TextEncodeShiftJIS)

	body, err = DecodeText(body, enc)
	// TODO: サイト変える
	fmt.Println(strings.Contains(string(body[:1000]), "総務省"))
	// Output:
	// <nil>
	// <nil>
	// true
	// true
}
