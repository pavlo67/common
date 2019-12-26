package importer_http_series

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/logger"

	"github.com/pavlo67/workshop/constructions/dataflow/flow_server_http"
	"github.com/pavlo67/workshop/constructions/dataimporter"
)

func NewSeriesHTTP(exportURL string, l logger.Operator) (dataimporter.Operator, error) {
	if l == nil {
		return nil, errors.New("on NewSeriesHTTP(): nil logger")
	}

	return &seriesHTTP{exportURL, l}, nil
}

var _ dataimporter.Operator = &seriesHTTP{}

type seriesHTTP struct {
	exportURL string
	l         logger.Operator
}

const onGet = "on seriesHTTP.Get(): "

func (sh *seriesHTTP) Get(lastImportedID string) (*dataimporter.DataSeries, error) {

	feedURL := sh.exportURL
	if lastImportedID != "" {
		feedURL += fmt.Sprintf("?%s=%s", flow_server_http.AfterIDParam, lastImportedID)
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

	series := dataimporter.DataSeries{}
	err = json.Unmarshal(body, &series)
	if err != nil {
		return nil, errors.Wrapf(err, onGet+"can't json.Unmarshal(%s, &series)", body)
	}

	return &series, nil
}
