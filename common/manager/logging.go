package manager

import (
	"bufio"
	"io"
	"os"

	"github.com/pavlo67/workshop/common/logger"
)

func Log(dataStream io.Reader, logfile string, l logger.Operator) {

	outfile, err := os.Create(logfile)
	if err != nil {
		l.Errorf("can't os.Create('%s')", logfile)

		// TODO??? panic(err)

		return
	}
	defer outfile.Close()

	io.Copy(bufio.NewWriter(outfile), dataStream)
}
