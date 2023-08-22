package config

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"runtime/debug"
	"strings"
	"testing"
	"time"

	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/logger/logger_zap"

	"github.com/stretchr/testify/require"
)

func ShowVCSInfo() {
	if bi, ok := debug.ReadBuildInfo(); ok && bi != nil {
		for _, s := range bi.Settings {
			if ok, _ := regexp.MatchString(`^vcs\.`, s.Key); ok {
				fmt.Printf("%s\t%s\n", s.Key, s.Value)
			}
		}
		fmt.Print("\n")
	}
}

func PrepareApp(envPath string) (Config, logger.Operator) {

	rand.Seed(time.Now().UnixNano())

	// get config --------------------------------------------------------------------

	configEnv, ok := os.LookupEnv("ENV")
	if !ok {
		configEnv = "local"
	}

	cfgServicePtr, err := Get(envPath+configEnv+".yaml", MarshalerYAML)
	if err != nil || cfgServicePtr == nil {
		log.Fatalf("on config.PrepareApp(%s, %s) got %#v / %s", envPath, configEnv+".yaml", cfgServicePtr, err)
	}

	// get logger --------------------------------------------------------------------

	// TODO!!! why it doesn't work??? it should be completed
	// var loggerConfig logger.Config
	// if err = cfgServicePtr.Value("logger", &loggerConfig); err != nil {
	//	log.Fatalf("on config.PrepareApp(%s, %s) got %s reading logger config", envPath, configEnv+".yaml", err)
	// }

	var loggerSaveFiles bool
	if err = cfgServicePtr.Value("logger_save_files", &loggerSaveFiles); err != nil {
		log.Fatalf("on config.PrepareApp(%s, %s) got %s reading 'logger_save_files'", envPath, configEnv+".yaml", err)
	}

	loggerConfig := logger.Config{
		SaveFiles: loggerSaveFiles,
	}

	l, err := logger_zap.New(loggerConfig)
	if err != nil {
		log.Fatal(err)
	}

	return *cfgServicePtr, l
}

func PrepareTests(t *testing.T, envPath, configEnv, logfile string) (Config, logger.Operator) {

	os.Setenv("ENV", configEnv)

	var logPath []string
	if logfile = strings.TrimSpace(logfile); logfile != "" {
		logPath = []string{logfile}
	}

	l, err := logger_zap.New(logger.Config{
		LogLevel:    logger.TraceLevel,
		OutputPaths: append(logPath, "stdout"),
		ErrorPaths:  append(logPath, "stderr"),
	}) // TODO!!! don't comment it (is required for tested components)
	require.NoError(t, err)
	require.NotNil(t, l)

	cfgServicePtr, err := Get(envPath+configEnv+".yaml", MarshalerYAML)
	require.NoError(t, err)
	require.NotNil(t, cfgServicePtr)

	return *cfgServicePtr, l

}
