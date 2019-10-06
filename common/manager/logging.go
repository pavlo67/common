package manager

import (
	"io"

	"github.com/pavlo67/workshop/common/logger"
)

func Redirect(key string, dataStream io.ReadCloser, outStream io.Writer, l logger.Operator) {
	defer dataStream.Close()
	if _, err := io.Copy(outStream, dataStream); err != nil {

		// TODO: insert key into each output line

		l.Errorf("%s: can't io.Copy(outStream, dataStream): %s", key, err)
	}
}
