package importer_series_http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/logger"

	"github.com/pavlo67/workshop/components/flow/flow_tagged/flow_tagged_server_http"
	"github.com/pavlo67/workshop/components/importer"
)

func NewSeriesHTTP(lastImportedID uint64, l logger.Operator) (importer.Operator, error) {
	if l == nil {
		return nil, errors.New("on NewSeriesHTTP(): nil logger")
	}

	return &seriesHTTP{lastImportedID, l}, nil
}

var _ importer.Operator = &seriesHTTP{}

type seriesHTTP struct {
	lastImportedID uint64
	l              logger.Operator
}

const onGet = "on seriesHTTP.Get(): "

func (sh *seriesHTTP) Get(feedURL string) (*importer.DataSeries, error) {
	if sh.lastImportedID > 0 {
		feedURL += fmt.Sprintf("?%s=%d", flow_tagged_server_http.AfterIDParam, sh.lastImportedID)
	}

	resp, err := http.Get(feedURL)
	if err != nil {
		return nil, errors.Wrapf(err, onGet+"can't http.Get(%s)", feedURL)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, onGet+"can't ioutil.ReadAll(%#v)", resp.Body)
	}

	series := importer.DataSeries{}
	err = json.Unmarshal(body, &series)
	if err != nil {
		return nil, errors.Wrapf(err, onGet+"can't json.Unmarshal(%s, &series)", body)
	}

	sh.lastImportedID = series.MaxID

	return &series, nil
}
