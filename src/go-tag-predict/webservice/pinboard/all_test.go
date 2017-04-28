package pinboard

import (
	"bytes"
	"fmt"
)

func ExampleParseAllData() {
	body := []byte(`
<posts>
	<post href="http://1" description="a" tag="">
	<post href="http://2" description="b" tag="tag">
	<post href="https://3" description="c" tag="にほんご b">
</posts>
	`)
	itr, err := parseAllData(bytes.NewReader(body))
	fmt.Println(err)
	i := 0
	for {
		post, err := itr.Next()
		if err != nil {
			fmt.Println(err)
			return
		}
		if post == nil {
			break
		}
		fmt.Println(post.Title)
		fmt.Println(post.Href)
		fmt.Println(post.Tags)
		i += 1
	}
	fmt.Println(i)

	// Output:
	// <nil>
	// a
	// http://1
	// []
	// b
	// http://2
	// [tag]
	// c
	// https://3
	// [にほんご b]
	// 3
}
