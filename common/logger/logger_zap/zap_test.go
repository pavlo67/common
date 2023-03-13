package logger_zap

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/logger"
)

func TestLoggerZapOutputPaths(t *testing.T) {

	cfg := logger.Config{
		LogLevel:    logger.InfoLevel,
		OutputPaths: []string{"test.log", "stdout"},
		//ErrorOutputPaths: nil,
		//Encoding:         "",
	}
	l, err := New(cfg)
	require.NoError(t, err)
	require.NotNil(t, l)

	l.Comment("INITIAL BLOCK")

	l.Info("1!!!")

	l.Comment("NEXT BLOCK")

	l.Error("2!!!")
}
