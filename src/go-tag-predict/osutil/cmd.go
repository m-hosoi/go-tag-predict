package osutil

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

// ExecuteCommand : 入出力付きで、外部プログラムを実行する
// $ echo input | name > []string
// ※ベタに実装するとstdinが詰まる現象が発生したので
func ExecuteCommand(ctx context.Context, name string, args []string, input string) ([]string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	var err error
	var stdin io.WriteCloser
	if input != "" {
		stdin, err = cmd.StdinPipe()
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	wg := sync.WaitGroup{}

	if input != "" {
		wg.Add(1)
		go func() {
			defer func() {
				stdin.Close()
				wg.Done()
			}()
			io.WriteString(stdin, input)
			io.WriteString(stdin, "\n")
		}()
	}

	err = cmd.Start()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// wg.Add(1)
	// go func() {
	// cmd.Wait()
	// wg.Done()
	// }()

	res := make([]string, 0, 2048)
	wg.Add(1)

	var readErr error
	go func() {
		out := bufio.NewReader(stdout)
		defer func() {
			stdout.Close()
			wg.Done()
		}()
		for {
			line, _, err := out.ReadLine()
			if err != nil {
				if err == io.EOF {
					break
				}
				if strings.Contains(err.Error(), "file already closed") {
					break
				}
				if strings.Contains(err.Error(), "bad file descriptor") {
					break
				}
				// fmt.Println("-------------------=")
				// fmt.Println(s)
				// fmt.Println(res)
				// fmt.Println("----")
				// fmt.Println(line)
				// fmt.Println("--")
				// panic(err)
				readErr = errors.WithStack(err)
				return
			}
			res = append(res, string(line))
		}
	}()
	wg.Wait()
	cmd.Wait()
	if readErr != nil {
		return nil, readErr
	}
	return res, nil
}
