package asyncwriter

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

func ExampleWriter() {
	ctx := context.Background()
	buf := bytes.Buffer{}
	aw := NewWriter(ctx, bufio.NewWriter(&buf), 16)

	wg := sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i != 99 {
				aw.WriteString(strconv.Itoa(i) + "\n")
			} else {
				aw.WriteString(strconv.Itoa(i))
			}
		}(i)
	}

	wg.Wait()
	aw.Close()
	<-aw.Done()
	fmt.Println(len(strings.Split(buf.String(), "\n")))

	// Output:
	// 100
}
