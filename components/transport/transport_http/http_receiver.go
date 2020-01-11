package transport_http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/server"
	"github.com/pavlo67/workshop/common/server/server_http"

	"github.com/pavlo67/workshop/components/packs"
)

func (transpOp *transportHTTP) receiveEndpoint() *server_http.Endpoint {
	if transpOp == nil {
		return nil
	}

	return &server_http.Endpoint{
		Method: "POST",
		WorkerHTTP: func(_ *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {
			var inPack packs.Pack

			packJSON, err := ioutil.ReadAll(req.Body)
			if err != nil {
				return server.ResponseRESTError(http.StatusBadRequest, errors.Errorf("ERROR on POST ...ReceivePack: reading body: %s", err))
			}

			err = json.Unmarshal(packJSON, &inPack)
			if err != nil {
				return server.ResponseRESTError(http.StatusBadRequest, errors.Errorf("ERROR on POST ...ReceivePack: can't json.Unmarshal(%s): %s", packJSON, err))
			}

			idIn, err := transpOp.packsOp.Save(&inPack, nil)
			if err != nil {
				return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on POST ...ReceivePack (%#v): '%s'", inPack, errors.Wrap(err, "can't transpOp.packsOp.Save(&inPack, nil)")))
			}

			runnerOp, taskID, err := transpOp.runnerFactory.TaskRunner(
				inPack.Task,
				&crud.SaveOptions{History: crud.History{{Key: crud.ProducedAction, DoneAt: time.Now(), RelatedIDs: []common.ID{idIn}}}},
				transpOp,
				inPack.From)
			if err != nil {
				return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on POST ...ReceivePack (%#v): '%s'", inPack, errors.Wrap(err, "can't transpOp.runnerFactory.TaskRunner(inPack.Task)")))
			}

			if taskID != "" {
				_, err = transpOp.packsOp.AddHistory(idIn, crud.History{{
					Key:        packs.TaskAction,
					DoneAt:     time.Now(),
					RelatedIDs: []common.ID{taskID},
					Errors:     nil,
				}}, nil)
				if err != nil {
					l.Error(err) // TODO!!! wrap the error
				}
			}

			_, err = runnerOp.Run()
			if err != nil {
				return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on POST ...ReceivePack (%#v): '%s'", inPack, errors.Wrap(err, "can't runnerOp.Run()")))
			}

			return server.ResponseRESTOk(common.Map{"target_task_id": taskID})
		},
	}
}
