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
	"github.com/pavlo67/workshop/common/identity"
	"github.com/pavlo67/workshop/common/server"
	"github.com/pavlo67/workshop/common/server/server_http"

	"github.com/pavlo67/workshop/components/packs"
)

const onAddHandler = "on transportHTTP.AddHandler(): "

func (transpOp *transportHTTP) AddHandler(receiverKey, typeKey identity.Key, handler packs.Handler) error {
	if handler == nil {
		return errors.Errorf(onAddHandler+"nil handler for key '%s'", typeKey)
	}

	transpOp.mutex.Lock()
	defer transpOp.mutex.Unlock()

	if _, ok := transpOp.handlers[[2]identity.Key{receiverKey, typeKey}]; ok {
		return errors.Errorf(onAddHandler+"A handler for key '%s' already exists", typeKey)
	}

	transpOp.handlers[[2]identity.Key{receiverKey, typeKey}] = handler
	return nil
}

const onRemoveHandler = "on transportHTTP.RemoveHandler(): "

func (transpOp *transportHTTP) RemoveHandler(receiverKey, typeKey identity.Key) {
	transpOp.mutex.Lock()
	defer transpOp.mutex.Unlock()

	if _, ok := transpOp.handlers[[2]identity.Key{receiverKey, typeKey}]; ok {
		delete(transpOp.handlers, [2]identity.Key{receiverKey, typeKey})
	}
}

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

			// TODO!!! ignore the domain name from .To

			transpOp.mutex.RLock()
			handler := transpOp.handlers[[2]identity.Key{inPack.To, inPack.TypeKey}]
			transpOp.mutex.RUnlock()
			if handler == nil {
				return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on POST ...ReceivePack (%#v): no handler for key '%s'", inPack, inPack.TypeKey))
			}

			idIn, err := transpOp.packsOp.Save(&inPack, nil)
			if err != nil {
				// TODO: wrap the error
				l.Error(err)
			}

			now := time.Now()

			packOut, err := handler.Handle(&inPack)
			if err != nil {
				return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on POST ...ReceivePack (%#v): '%s'", inPack, err))
			}
			packOut.History = append(packOut.History, crud.Action{
				Key:        crud.ProducedAction,
				DoneAt:     now,
				RelatedIDs: []common.ID{idIn},
			})
			var producedIDs []common.ID
			if packOut != nil {
				idOut, err := transpOp.packsOp.Save(packOut, nil)
				if err != nil {
					// TODO: wrap the error
					l.Error(err)
				} else {
					producedIDs = []common.ID{idOut}
				}
			}

			_, err = transpOp.packsOp.AddHistory(idIn, crud.History{{
				Key:        packs.HandleAction,
				DoneAt:     now,
				RelatedIDs: producedIDs,
			}}, nil)
			if err != nil {
				// TODO: wrap the error
				l.Error(err)
			}

			return server.ResponseRESTOk(packOut)
		},
	}
}
