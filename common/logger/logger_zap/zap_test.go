package logger_zap

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/logger"
)

func TestLoggerZapOutputPaths(t *testing.T) {

	testKey := "testKey"

	cfg := logger.Config{
		Key:         testKey,
		LogLevel:    logger.InfoLevel,
		OutputPaths: []string{"test.log", "stdout"},
		//ErrorOutputPaths: nil,
		//Encoding:         "",
	}
	l, err := New(cfg)
	require.NoError(t, err)
	require.NotNil(t, l)

	require.Equal(t, testKey, l.Key())

	l.Comment("INITIAL BLOCK")

	l.Info("IT'S TEST FOR LOGGING INFO. OK!")

	l.Comment("NEXT BLOCK")

	l.Error("IT'S TEST FOR LOGGING ERROR. OK!")
}
