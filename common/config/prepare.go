package config

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strconv"
	"testing"
	"time"

	"github.com/pavlo67/common/common/colorize"

	"github.com/pavlo67/common/common/serialization"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/filelib"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/logger/logger_zap"
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

func PrepareApp(envPath, logPath string) (Envs, logger.Operator) {

	rand.Seed(time.Now().UnixNano())

	// get config --------------------------------------------------------------------

	configEnv, ok := os.LookupEnv("ENV")
	if !ok {
		configEnv = "local"
	}

	cfgServicePtr, err := Get(envPath+configEnv+".yaml", serialization.MarshalerYAML)
	if err != nil || cfgServicePtr == nil {
		log.Fatalf("on config.PrepareApp(%s, %s) got %#v / %s", envPath, configEnv+".yaml", cfgServicePtr, err)
	}

	// get logger --------------------------------------------------------------------

	var saveFiles bool
	if err = cfgServicePtr.Value("logger_save_files", &saveFiles); err != nil {
		fmt.Fprintf(os.Stderr, colorize.Red+"on config.PrepareApp(%s, %s), reading of 'logger_save_files' key produces the error: %s\n"+colorize.Reset, envPath, configEnv+".yaml", err)
	}

	if logPath == "" {
		if err = cfgServicePtr.Value("logger_path", &logPath); err != nil {
			fmt.Fprintf(os.Stderr, colorize.Red+"on config.PrepareApp(%s, %s), reading of 'logger_path' key produces the error: %s\n"+colorize.Reset, envPath, configEnv+".yaml", err)
		}
	}

	logPath, err = filelib.Dir(logPath)
	if err != nil {
		log.Fatalf("on config.PrepareApp(%s, %s): can't create log path (%s): %s", envPath, configEnv+".yaml", logPath, err)
	}

	loggerConfig := logger.Config{
		Key:       strconv.FormatInt(time.Now().Unix(), 10),
		BasePath:  logPath,
		SaveFiles: saveFiles,
	}

	l, err := logger_zap.New(loggerConfig)
	if err != nil {
		log.Fatal(err)
	}

	return *cfgServicePtr, l
}

func PrepareTests(t *testing.T, envPath, logFile string) (Envs, logger.Operator) {

	configEnv := "test"
	os.Setenv("ENV", configEnv)

	envsPtr, err := Get(envPath+configEnv+".yaml", serialization.MarshalerYAML)
	require.NoError(t, err)
	require.NotNil(t, envsPtr)

	var loggerFiles []string
	if logFile != "" {
		loggerFiles = append(loggerFiles, logFile)
	}

	var loggerSaveFiles bool
	if err = envsPtr.Value("logger_save_files", &loggerSaveFiles); err != nil {
		fmt.Fprintf(os.Stderr, colorize.Red+"on config.PrepareApp(%s, %s), reading of 'logger_save_files' key produces the error: %s\n"+colorize.Reset, envPath, configEnv+".yaml", err)
	}

	var loggerPath string
	if err = envsPtr.Value("logger_path", &loggerPath); err != nil {
		fmt.Fprintf(os.Stderr, colorize.Red+"on config.PrepareApp(%s, %s), reading of 'logger_path' key produces the error: %s\n"+colorize.Reset, envPath, configEnv+".yaml", err)
	}

	logKey := time.Now().Format(time.RFC3339)[:19]
	logPath, err := filelib.Dir(filepath.Join(loggerPath, logKey))
	require.NoError(t, err)

	l, err := logger_zap.New(logger.Config{
		Key:         logKey,
		BasePath:    logPath,
		SaveFiles:   loggerSaveFiles,
		LogLevel:    logger.TraceLevel,
		OutputPaths: append(loggerFiles, "stdout"),
		ErrorPaths:  append(loggerFiles, "stderr"),
	})
	require.NoError(t, err)
	require.NotNil(t, l)

	return *envsPtr, l

}

//func PrepareTests(t *testing.T, envPath, configEnv, logfile string) (Envs, logger.Operator) {
//
//	os.Setenv("ENV", configEnv)
//
//	var logPath []string
//	if logfile = strings.TrimSpace(logfile); logfile != "" {
//		logPath = []string{logfile}
//	}
//
//	l, err := logger_zap.New(logger.Config{
//		LogLevel:    logger.TraceLevel,
//		OutputPaths: append(logPath, "stdout"),
//		ErrorPaths:  append(logPath, "stderr"),
//	}) // TODO!!! don't comment it (is required for tested components)
//	require.NoError(t, err)
//	require.NotNil(t, l)
//
//	cfgServicePtr, err := Get(envPath+configEnv+".yaml", serialization.MarshalerYAML)
//	require.NoError(t, err)
//	require.NotNil(t, cfgServicePtr)
//
//	return *cfgServicePtr, l
//
//}
