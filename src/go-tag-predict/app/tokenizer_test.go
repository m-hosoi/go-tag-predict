package app

import (
	"context"
	"fmt"
	"go-tag-predict/fileutil"
	"path"
)

func ExampleTokenizeMecab() {
	configPath := path.Join(fileutil.GetGopath()[0], "go-tag-predict.toml")
	config, err := LoadConfig(configPath)

	ar, err := tokenizeMecab(context.Background(), config.Mecab, "")
	fmt.Println(err)
	fmt.Println(len(ar))

	ar, err = tokenizeMecab(context.Background(), config.Mecab, "あのイーハトーヴォのすきとおった風、夏でも底に冷たさをもつ青いそら、\nうつくしい森で飾られたモリーオ市、郊外のぎらぎらひかる草の波。")
	fmt.Println(err)
	for _, s := range ar {
		fmt.Println(s)
	}
	// Output:
	// <nil>
	// 0
	// <nil>
	// あの
	// イーハトーヴォ
	// の
	// すきとおっ
	// た
	// 風
	// 、
	// 夏
	// で
	// も
	// 底
	// に
	// 冷た
	// さ
	// を
	// もつ
	// 青い
	// そら
	// 、
	// うつくしい
	// 森
	// で
	// 飾ら
	// れ
	// た
	// モリーオ
	// 市
	// 、
	// 郊外
	// の
	// ぎらぎら
	// ひかる
	// 草
	// の
	// 波
	// 。
}
func ExampleTokenizeJumanpp() {
	configPath := path.Join(fileutil.GetGopath()[0], "go-tag-predict.toml")
	config, err := LoadConfig(configPath)

	ar, err := tokenizeMecab(context.Background(), config.Mecab, "")
	fmt.Println(err)
	fmt.Println(len(ar))

	ar, err = tokenizeJumanpp(context.Background(), config.Jumanpp, "あのイーハトーヴォのすきとおった風、夏でも底に冷たさをもつ青いそら、\nうつくしい森で飾られたモリーオ市、郊外のぎらぎらひかる草の波。")
	fmt.Println(err)
	for _, s := range ar {
		fmt.Println(s)
	}
	// Output:
	// <nil>
	// 0
	// <nil>
	// あの
	// イーハトーヴォ
	// の
	// すきとおった
	// 風
	// 、
	// 夏
	// でも
	// 底
	// に
	// 冷た
	// さ
	// を
	// もつ
	// 青い
	// そら
	// 、
	// うつくしい
	// 森
	// で
	// 飾ら
	// れた
	// モリーオ
	// 市
	// 、
	// 郊外
	// の
	// ぎらぎら
	// ひかる
	// 草
	// の
	// 波
	// 。
}
