package main

import (
	"context"
	"flag"
	"fmt"
	"go-tag-predict/app"
	"go-tag-predict/fileutil"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"

	"net/http"
	_ "net/http/pprof"

	"go.uber.org/zap"
)

var configPath *string
var isDebugMode = flag.Bool("debug", false, "debug mode")

func init() {
	configPath = flag.String("f", fileutil.FindFilePath("go-tag-predict.toml"), "configuration file name")
}

func main() {
	logger, err := createLogger()
	if err != nil {
		checkErrorExit(err)
	}

	flag.Usage = func() {
		fmt.Printf("USAGE: %s [options] COMMAND\n\n", filepath.Base(os.Args[0]))
		fmt.Printf("Commands:\n")
		fmt.Printf("  supervised  学習モード\n")
		fmt.Printf("  predict     分類モード\n")
		fmt.Printf("  help        Print this message\n")
		fmt.Printf("\n")
		fmt.Printf("Run '%s COMMAND --help' for more information on the command\n", filepath.Base(os.Args[0]))
		fmt.Printf("\n")
		fmt.Printf("Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	if *isDebugMode {
		logger.Info("debug mode")
		go func() {
			run(context.Background(), logger)
		}()
		// http://localhost:6060/debug/pprof/
		// http://localhost:6060/debug/pprof/goroutine?debug=1
		log.Println(http.ListenAndServe("localhost:6060", nil))
	} else {
		run(context.Background(), logger)
	}
}
func run(ctx context.Context, logger *zap.Logger) {
	logger.Info("start", zap.Int("numCPU", runtime.NumCPU()), zap.Int("maxProcs", runtime.GOMAXPROCS(0)))
	config, err := app.LoadConfig(*configPath)
	checkErrorExit(err)

	command := ""
	// command := "predict" // debug

	args := flag.Args()
	if len(args) >= 1 {
		command = args[0]
	}
	switch command {
	case "supervised":
		err = app.RunSupervised(ctx, config, logger)
		checkErrorExit(err)
	case "predict":
		err = app.RunPredict(ctx, config, logger)
		checkErrorExit(err)
	case "help":
		flag.Usage()
	default:
		fmt.Printf("%q is not valid command (supervised or predict).\n\n", command)
		flag.Usage()
		os.Exit(1)
	}
	logger.Info("finish")
}
func checkErrorExit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
func createLogger() (*zap.Logger, error) {
	var logger *zap.Logger
	var err error
	if *isDebugMode {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return logger, nil
}
