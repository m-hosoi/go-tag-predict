package app

import (
	"context"
	"go-tag-predict/lambda"
	"go-tag-predict/osutil"
	"strings"

	"github.com/pkg/errors"
	mecab "github.com/shogo82148/go-mecab"
)

// Tokenize : テキストを単語で分割する
func Tokenize(ctx context.Context, config *Config, s string) ([]string, error) {
	if config.Tokenizer == "mecab" {
		return tokenizeMecab(ctx, config.Mecab, s)
	} else if config.Tokenizer == "jumanpp" {
		return tokenizeJumanpp(ctx, config.Jumanpp, s)
	}
	panic("bad tokenizer (mecab or jumanpp)")
}
func tokenizeMecab(ctx context.Context, config *MecabConfig, s string) ([]string, error) {
	tagger, err := mecab.New(map[string]string{
		"output-format-type": "wakati",
		"dicdir":             config.DictDirPath,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer tagger.Destroy()

	// 入力データが大きすぎると、Mecabのエラー「too long sentence.」が発生する
	if len(s) > 1024*256 {
		s = s[:1024*256]
	}
	res, err := tagger.Parse(s)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	lines := strings.Split(res, " ")
	lines = lambda.MapString(lines, func(s string) string {
		return strings.TrimSpace(s)
	})
	lines = lambda.FilterString(lines, func(s string) bool {
		return s != ""
	})
	return lines, nil
}
func tokenizeJumanpp(ctx context.Context, config *JumanppConfig, s string) ([]string, error) {
	lines, err := osutil.ExecuteCommand(ctx, config.Command, config.Args, s)
	if err != nil {
		return nil, err
	}
	lines = lambda.MapString(lines, func(s string) string {
		ar := strings.SplitN(string(s), config.TokenSeparator, 2)
		if len(ar) != 2 {
			return ""
		}
		return ar[0]
	})
	lines = lambda.FilterString(lines, func(s string) bool {
		return s != ""
	})
	return lines, nil
}
