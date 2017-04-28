package asyncwriter

import (
	"bufio"
	"context"
)

// Writer : 複数Goルーチンから非同期出力可能なWriter
// Example:
//   aw := newAsyncWriter(ctx, bufio.NewWriter(f))
//   go func() {
//     ...
//     aw.WriteString(...)
//   }()
//   ...
//   aw.Close()
//   <-aw.Done()
type Writer struct {
	ch     chan string
	done   chan struct{}
	writer *bufio.Writer
}

// NewWriter : コンストラクタ
func NewWriter(ctx context.Context, w *bufio.Writer, bufferCount int) *Writer {
	aw := &Writer{
		ch:     make(chan string, bufferCount),
		done:   make(chan struct{}, 0),
		writer: w,
	}
	aw.start(ctx)
	return aw
}
func (a *Writer) start(ctx context.Context) {
	go func() {
	FOR:
		for {
			select {
			case <-ctx.Done():
				break FOR
			case s, ok := <-a.ch:
				if !ok {
					break FOR
				}
				a.writer.WriteString(s)
			}
		}
		a.writer.Flush()
		a.done <- struct{}{}
	}()
}

// WriteString : stringを出力する
func (a *Writer) WriteString(s string) {
	a.ch <- s
}

// Close : 書き込み完了
func (a *Writer) Close() {
	close(a.ch)
}

// Done : Close完了まで待機するためのchannelを取得する
func (a *Writer) Done() chan struct{} {
	return a.done
}
