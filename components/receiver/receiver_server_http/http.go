package receiver_server_http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/server"
	"github.com/pavlo67/workshop/common/server/server_http"
	"github.com/pavlo67/workshop/common/types"

	"github.com/pavlo67/workshop/components/packs"
	"github.com/pavlo67/workshop/components/receiver"
)

var _ receiver.Operator = &receiverHTTP{}

type receiverHTTP struct {
	handlers map[types.Key]packs.Handler

	packsOp packs.Operator
	mutex   *sync.RWMutex
}

const onNew = "on receiverHTTP.New(): "

func New(packsOp packs.Operator) (receiver.Operator, *server_http.Endpoint, error) {
	if packsOp == nil {
		return nil, nil, errors.New(onNew + "no packs.Operator")
	}

	handlers := map[types.Key]packs.Handler{}

	receiverOp := &receiverHTTP{
		handlers: handlers,
		packsOp:  packsOp,
		mutex:    &sync.RWMutex{},
	}
	return receiverOp, receiverOp.receiveEndpoint(), nil
}

const onAddHandler = "on receiverHTTP.AddHandler(): "

func (receiverOp *receiverHTTP) AddHandler(typeKey types.Key, handler packs.Handler) error {
	if handler == nil {
		return errors.Errorf(onAddHandler+"nil handler for key '%s'", typeKey)
	}

	receiverOp.mutex.Lock()
	defer receiverOp.mutex.Unlock()

	if _, ok := receiverOp.handlers[typeKey]; ok {
		return errors.Errorf(onAddHandler+"A handler for key '%s' already exists", typeKey)
	}

	receiverOp.handlers[typeKey] = handler
	return nil
}

const onRemoveHandler = "on receiverHTTP.RemoveHandler(): "

func (receiverOp *receiverHTTP) RemoveHandler(typeKey types.Key) {
	receiverOp.mutex.Lock()
	defer receiverOp.mutex.Unlock()

	if _, ok := receiverOp.handlers[typeKey]; ok {
		delete(receiverOp.handlers, typeKey)
	}
}

func (receiverOp *receiverHTTP) receiveEndpoint() *server_http.Endpoint {
	if receiverOp == nil {
		return nil
	}

	return &server_http.Endpoint{
		Method: "POST",
		WorkerHTTP: func(_ *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {
			var packIn packs.Pack

			packJSON, err := ioutil.ReadAll(req.Body)
			if err != nil {
				return server.ResponseRESTError(http.StatusBadRequest, errors.Errorf("ERROR on POST ...ReceivePack: reading body: %s", err))
			}

			err = json.Unmarshal(packJSON, &packIn)
			if err != nil {
				return server.ResponseRESTError(http.StatusBadRequest, errors.Errorf("ERROR on POST ...ReceivePack: can't json.Unmarshal(%s): %s", packJSON, err))
			}

			receiverOp.mutex.RLock()
			handler := receiverOp.handlers[packIn.TypeKey]
			receiverOp.mutex.RUnlock()
			if handler == nil {
				return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on POST ...ReceivePack (%#v): no handler for key '%s'", packIn, packIn.TypeKey))
			}

			idIn, err := receiverOp.packsOp.Save(&packIn, nil)
			if err != nil {
				// TODO: wrap the error
				l.Error(err)
			}

			now := time.Now()

			packOut, err := handler.Handle(&packIn)
			if err != nil {
				return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on POST ...ReceivePack (%#v): '%s'", packIn, err))
			}
			packOut.History = append(packOut.History, crud.Action{
				Key:        crud.ProducedAction,
				DoneAt:     now,
				RelatedIDs: []common.ID{idIn},
			})
			var producedIDs []common.ID
			if packOut != nil {
				idOut, err := receiverOp.packsOp.Save(packOut, nil)
				if err != nil {
					// TODO: wrap the error
					l.Error(err)
				} else {
					producedIDs = []common.ID{idOut}
				}
			}

			_, err = receiverOp.packsOp.AddHistory(idIn, crud.History{{
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
