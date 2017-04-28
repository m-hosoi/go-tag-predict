package app

import (
	"go-tag-predict/fileutil"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

// Config : 設定ファイル
type Config struct {
	CacheDirPath string `toml:"cache_dir"`
	TmpDirPath   string `toml:"tmp_dir"`
	Tokenizer    string `toml:"tokenizer"`
	Supervised   *SupervisedConfig
	Predict      *PredictConfig
	Fasttext     *FasttextConfig
	Mecab        *MecabConfig
	Jumanpp      *JumanppConfig
}

// PredictConfig : 学習処理の設定
type SupervisedConfig struct {
	LearningSourceFilePath string `toml:"learning_source_file"`
	ParallelsCount         int    `toml:"parallels_count"`
	WriterBufferSize       int    `toml:"writer_buffer_size"`
	WriterQueueCount       int    `toml:"writer_queue_count"`
}

// PredictConfig : 分類処理の設定
type PredictConfig struct {
	FeedURLs       []string `toml:"feed_urls"`
	ParallelsCount int      `toml:"parallels_count"`
	MinProbability float64  `toml:"min_probability"`
}

// FasttextConfig : fastTextの設定
type FasttextConfig struct {
	Command        string   `toml:"command"`
	SupervisedArgs []string `toml:"supervised_args"`
	PredictArgs    []string `toml:"predict_args"`
}

// MecabConfig : Mecabの設定
type MecabConfig struct {
	DictDirPath string `toml:"dict_dir"`
}

// JumanppConfig : juman++の設定
type JumanppConfig struct {
	Command        string   `toml:"command"`
	Args           []string `toml:"args"`
	TokenSeparator string   `toml:"token_separator"`
}

// NewConfig : Configのコンストラクタ
func NewConfig() *Config {
	return &Config{}
}

// LoadConfig : Configをtomlの設定ファイルからロード
func LoadConfig(configPath string) (*Config, error) {
	config := NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	config.Supervised.LearningSourceFilePath = fileutil.FindFilePath(config.Supervised.LearningSourceFilePath)
	config.CacheDirPath = fileutil.FindFilePath(config.CacheDirPath)
	config.TmpDirPath = fileutil.FindFilePath(config.TmpDirPath)

	return config, nil
}

// GetModelPathForSupervised : fasttext学習データの場所(supervised用の形式)
func (c *Config) GetModelPathForSupervised() string {
	return path.Join(c.TmpDirPath, "model")
}

// GetModelPathForPredict : fasttext学習データの場所(predict用の形式)
func (c *Config) GetModelPathForPredict() string {
	return c.GetModelPathForSupervised() + ".bin"
}

// GetTagIDPath : 単語=>数値変換表の場所
func (c *Config) GetTagIDPath() string {
	return path.Join(c.TmpDirPath, "tagid.txt")
}

// GetSupervisedSourcePath : fasttext学習の入力データの場所
func (c *Config) GetSupervisedSourcePath() string {
	return path.Join(c.TmpDirPath, "input.txt")
}
