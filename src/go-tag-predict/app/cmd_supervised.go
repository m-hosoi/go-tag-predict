package app

import (
	"bufio"
	"context"
	"go-tag-predict/asyncwriter"
	"go-tag-predict/lambda"
	"go-tag-predict/webservice/pinboard"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/pkg/errors"

	"golang.org/x/sync/errgroup"
)

// RunSupervised : 学習メイン関数
func RunSupervised(ctx context.Context, config *Config, logger *zap.Logger) error {
	err := createSupervisedInput(ctx, config, logger)
	if err != nil {
		return err
	}
	err = supervised(ctx, config)
	if err != nil {
		return err
	}
	return nil
}
func createSupervisedInput(ctx context.Context, config *Config, logger *zap.Logger) error {
	itr, err := pinboard.LoadFile(config.Supervised.LearningSourceFilePath)
	if err != nil {
		return err
	}
	err = os.MkdirAll(config.TmpDirPath, 0700)
	if err != nil {
		return errors.WithStack(err)
	}

	f, err := os.Create(config.GetSupervisedSourcePath())
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()

	aw := asyncwriter.NewWriter(ctx, bufio.NewWriterSize(f, config.Supervised.WriterBufferSize), config.Supervised.WriterQueueCount)

	t := NewTagID()

	eg, ctx := errgroup.WithContext(ctx)
	limitter := make(chan struct{}, max(0, config.Supervised.ParallelsCount-1)) // 同時実行数の制御
	i := 0
	for {
		post, err := itr.Next()
		if err != nil {
			return err
		}
		if post == nil {
			break
		}
		limitter <- struct{}{}
		if ctx.Err() != nil {
			break
		}
		func(post *pinboard.Post) {
			eg.Go(func() error {
				defer func() {
					<-limitter
				}()
				return procPost(ctx, config, logger, t, aw, post)
			})
		}(post)
		i++
		// // // DEBUG
		// if i == 200 {
		// fmt.Println("DEBUG")
		// }
		// if i > 200 {
		// time.Sleep(60 * time.Second)
		// }
	}
	err = eg.Wait()
	if err != nil {
		return err
	}
	// fmt.Println(i)

	aw.Close()
	<-aw.Done()

	if f, err := os.Create(config.GetTagIDPath()); err == nil {
		defer f.Close()
		t.Serialize(f)
	} else {
		return errors.WithStack(err)
	}
	return nil
}
func supervised(ctx context.Context, config *Config) error {
	modelPath := config.GetModelPathForSupervised()
	inputPath := config.GetSupervisedSourcePath()

	cmd := exec.CommandContext(
		ctx,
		config.Fasttext.Command,
		lambda.MapString(config.Fasttext.SupervisedArgs, func(s string) string {
			if s == "{MODEL_PATH}" {
				return modelPath
			}
			if s == "{DATA_PATH}" {
				return inputPath
			}
			return s
		})...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
func procPost(ctx context.Context, config *Config, logger *zap.Logger, t TagID, aw *asyncwriter.Writer, post *pinboard.Post) error {
	logger.Debug("begin",
		zap.String("url", post.Href),
		zap.Int("goroutines", runtime.NumGoroutine()),
	)
	// time.Sleep(600 * time.Second)
	content, err := LoadWebContent(ctx, post.Href, config.CacheDirPath)
	logger.Debug("get",
		zap.String("url", post.Href),
		zap.Int("goroutines", runtime.NumGoroutine()),
	)
	if err != nil {
		logger.Debug("skip",
			zap.String("url", post.Href),
			zap.String("err", err.Error()),
			zap.Int("goroutines", runtime.NumGoroutine()),
		)
		return nil // ページの取得に失敗しても全体の処理を継続する
	}
	if len(content) < 128 { // 本文が短いデータを除去
		return nil
	}
	tokens, err := tokenizeWebContent(ctx, config, post, content)
	logger.Debug("tokenize",
		zap.String("url", post.Href),
		zap.Int("goroutines", runtime.NumGoroutine()),
	)
	if err != nil {
		logger.Debug("tokenize error",
			zap.String("url", post.Href),
			zap.Int("goroutines", runtime.NumGoroutine()),
		)
		return err
	}
	// TODO: 句読点などのノイズを除去

	tokens = lambda.MapIntString(t.GetIDs(tokens), strconv.Itoa)

	strTokens := strings.Join(tokens, " ")
	for _, tag := range post.Tags {
		aw.WriteString(
			strings.Join([]string{"__label__", strconv.Itoa(t.GetID(tag)), " , ", strTokens, "\n"}, ""),
		)
	}
	logger.Info("finish",
		zap.String("url", post.Href),
		zap.Int("goroutines", runtime.NumGoroutine()),
	)
	return nil
}
func tokenizeWebContent(ctx context.Context, config *Config, post *pinboard.Post, content string) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	type resultSet struct {
		v   []string
		err error
	}
	c := make(chan resultSet, 0)
	go func() {
		if v, err := Tokenize(ctx, config, post.Title+"\n"+content); err != nil {
			c <- resultSet{nil, errors.Wrap(err, "tokenize error: \n"+post.Title+"\n"+post.Href+"\nSize:"+strconv.Itoa(len(content)))}
		} else {
			c <- resultSet{v, nil}
		}
	}()
	select {
	case <-ctx.Done():
		return nil, errors.WithStack(ctx.Err())
	case r := <-c:
		if r.err != nil {
			return nil, r.err
		}
		return r.v, nil
	}
}
func max(i0 int, i1 int) int {
	if i0 > i1 {
		return i0
	}
	return i1
}
