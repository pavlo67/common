package apps

import (
	"log"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/pavlo67/common/common/logger/logger_zap"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/logger"
)

func Prepare(envsSubpath string) (envPath string, cfgService *config.Config, l logger.Operator) {

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

	cwd, err := os.Getwd()
	if err != nil {
		l.Fatal("can't os.Getwd(): ", err)
	}
	cwd += "/"

	envPath = cwd + envsSubpath
	cfgServicePath := envPath + configEnv + ".yaml"
	cfgService, err = config.Get(cfgServicePath, config.MarshalerYAML)
	if err != nil || cfgService == nil {
		l.Fatalf("on config.Get(%s, serializer.MarshalerYAML)", cfgServicePath, cfgService, err)
	}

	return envPath, cfgService, l
}

func PrepareTests(t *testing.T, envsSubpath, configEnv, logfile string) (envPath string, cfgService *config.Config, l logger.Operator) {

	os.Setenv("ENV", configEnv)

	var logPath []string
	if logfile = strings.TrimSpace(logfile); logfile != "" {
		logPath = []string{logfile}
	}

	l, err := logger_zap.New(logger.Config{
		LogLevel:         logger.TraceLevel,
		OutputPaths:      append(logPath, "stdout"),
		ErrorOutputPaths: append(logPath, "stderr"),
		Encoding:         "",
	}) // TODO!!! don't comment it (is required for tested components)
	require.NoError(t, err)
	require.NotNil(t, l)

	cwd, err := os.Getwd()
	require.NoError(t, err)

	cwd += "/"

	envPath = cwd + envsSubpath
	cfgServicePath := envPath + configEnv + ".yaml"
	cfgService, err = config.Get(cfgServicePath, config.MarshalerYAML)
	require.NoError(t, err)
	require.NotNil(t, cfgService)

	return envPath, cfgService, l

}
