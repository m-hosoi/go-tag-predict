package app

import (
	"context"
	"go-tag-predict/lambda"
	"go-tag-predict/osutil"
	"os"
	"strconv"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// RunPredict : 分類メイン関数
func RunPredict(ctx context.Context, config *Config, logger *zap.Logger) error {
	f, err := os.Open(config.GetTagIDPath())
	if err != nil {
		return errors.WithStack(err)
	}
	t, err := LoadTagID(f)
	if err != nil {
		return err
	}
	idMap := t.GetReverse()

	eg, ctx := errgroup.WithContext(ctx)
	limitter := make(chan struct{}, max(0, config.Predict.ParallelsCount-1)) // 同時実行数の制御
	for _, rawurl := range config.Predict.FeedURLs {
		feed, err := LoadFeed(ctx, rawurl, config.CacheDirPath)
		if err != nil {
			return err
		}
		for _, item := range feed.Items {
			limitter <- struct{}{}
			if ctx.Err() != nil {
				break
			}
			func(item *gofeed.Item) {
				eg.Go(func() error {
					defer func() {
						<-limitter
					}()
					return procPage(ctx, config, logger, t, idMap, item)
				})
			}(item)
		}
	}
	err = eg.Wait()
	if err != nil {
		return err
	}

	return nil
}
func procPage(ctx context.Context, config *Config, logger *zap.Logger, t TagID, idMap map[int]string, item *gofeed.Item) error {
	content, err := LoadWebContent(ctx, item.Link, config.CacheDirPath)
	if err != nil {
		return nil // ページの取得に失敗しても全体の処理を継続する
	}
	tokens, err := Tokenize(ctx, config, content)
	tokens = lambda.MapIntString(t.GetIDs(tokens), strconv.Itoa)
	if err != nil {
		return err
	}

	ar, err := predict(ctx, config, idMap, strings.Join(tokens, " "))
	if err != nil {
		return err
	}
	if len(ar) == 0 {
		return nil
	}
	logger.Info("page",
		zap.String("url", item.Link),
		zap.Array("tag", ar))
	return nil
}
func predict(ctx context.Context, config *Config, idMap map[int]string, tokens string) (predictResults, error) {
	modelPath := config.GetModelPathForPredict()
	lines, err := osutil.ExecuteCommand(
		ctx,
		config.Fasttext.Command,
		lambda.MapString(config.Fasttext.PredictArgs, func(s string) string {
			if s == "{MODEL_PATH}" {
				return modelPath
			}
			return s
		}), tokens)
	if err != nil {
		return nil, err
	}
	res := make(predictResults, 0, len(lines))
	for _, line := range lines {
		ar := strings.SplitN(line, " ", 2)
		id, err := strconv.Atoi(ar[0][len("__label__"):])
		if err != nil {
			return nil, errors.WithStack(err)
		}
		probability, err := strconv.ParseFloat(ar[1], 32)
		if probability < config.Predict.MinProbability {
			continue
		}
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if tag, ok := idMap[id]; ok {
			res = append(res, &predictResult{tag: tag, probability: probability})
		}
	}
	return res, nil
}

type predictResult struct {
	tag         string
	probability float64
}

type predictResults []*predictResult

func (ps predictResults) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, p := range ps {
		enc.AppendObject(p)
	}
	return nil
}

func (p *predictResult) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if p != nil {
		enc.AddString("tag", p.tag)
		enc.AddFloat64("probability", p.probability)
	}
	return nil
}
