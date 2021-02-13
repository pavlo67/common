package apps

import (
	"flag"
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

func Prepare(buildDate, buildTag, buildCommit, serviceName, appsSubpathDefault string) (versionOnly bool, envPath string, cfgService *config.Config, l logger.Operator) {

	rand.Seed(time.Now().UnixNano())

	// show build/console params -----------------------------------------------------

	var appsSubpath string
	flag.BoolVar(&versionOnly, "v", false, "show build vars only")
	flag.StringVar(&appsSubpath, "apps_subpath", appsSubpathDefault, "subpath to /apps directory")
	flag.Parse()

	if buildDate = strings.TrimSpace(buildDate); buildDate == "" {
		buildDate = time.Now().Format(time.RFC3339)
	}

	log.Printf("builded: %s, tag: %s, commit: %s\n", buildDate, buildTag, buildCommit)

	if versionOnly {
		return versionOnly, "", nil, nil
	}

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

	envPath = cwd + appsSubpath + "_environments/"
	cfgServicePath := envPath + configEnv + ".yaml"
	cfgService, err = config.Get(cfgServicePath, serviceName, config.MarshalerYAML)
	if err != nil || cfgService == nil {
		l.Fatalf("on config.Get(%s, %s, serializer.MarshalerYAML)", cfgServicePath, serviceName, cfgService, err)
	}

	return versionOnly, envPath, cfgService, l
}

func PrepareTests(t *testing.T, serviceName, appsSubpath, configEnv, logfile string) (envPath string, cfgService *config.Config, l logger.Operator) {

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

	envPath = cwd + appsSubpath + "_environments/"
	cfgServicePath := envPath + configEnv + ".yaml"
	cfgService, err = config.Get(cfgServicePath, serviceName, config.MarshalerYAML)
	require.NoError(t, err)
	require.NotNil(t, cfgService)

	return envPath, cfgService, l

}
