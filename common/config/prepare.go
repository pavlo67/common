package config

import (
	"log"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/logger/logger_zap"

	"github.com/stretchr/testify/require"
)

func PrepareApp(envPath string) (Config, logger.Operator) {

	rand.Seed(time.Now().UnixNano())

	// get logger --------------------------------------------------------------------

	l, err := logger_zap.New(logger.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// get config --------------------------------------------------------------------

	configEnv, ok := os.LookupEnv("ENV")
	if !ok {
		configEnv = "local"
	}

	cfgServicePtr, err := Get(envPath+configEnv+".yaml", MarshalerYAML)
	if err != nil || cfgServicePtr == nil {
		l.Fatalf("on config.PrepareApp(%s, %s) got %#v / %s", envPath, configEnv+".yaml", cfgServicePtr, err)
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
		LogLevel:         logger.TraceLevel,
		OutputPaths:      append(logPath, "stdout"),
		ErrorOutputPaths: append(logPath, "stderr"),
	}) // TODO!!! don't comment it (is required for tested components)
	require.NoError(t, err)
	require.NotNil(t, l)

	cfgServicePtr, err := Get(envPath+configEnv+".yaml", MarshalerYAML)
	require.NoError(t, err)
	require.NotNil(t, cfgServicePtr)

	return *cfgServicePtr, l

}
