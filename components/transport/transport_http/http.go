package transport_http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/identity"
	"github.com/pavlo67/workshop/common/server"
	"github.com/pavlo67/workshop/common/server/server_http"

	"github.com/pavlo67/workshop/components/packs"
	"github.com/pavlo67/workshop/components/transport"
	"github.com/pavlo67/workshop/components/transportrouter"
)

var _ transport.Operator = &transportHTTP{}

type transportHTTP struct {
	packsOp  packs.Operator
	routerOp transportrouter.Operator

	domain identity.Domain
	path   string
	id     uint64

	routes transportrouter.Routes

	handlers map[[2]identity.Key]packs.Handler
	mutex    *sync.RWMutex
}

const onNew = "on sender_http.New(): "

func New(packsOp packs.Operator, routerOp transportrouter.Operator, domain identity.Domain) (transport.Operator, *server_http.Endpoint, error) {
	if packsOp == nil {
		return nil, nil, errors.New(onNew + "no packs.Operator")
	}

	if routerOp == nil {
		return nil, nil, errors.New(onNew + "no router.Operator")
	}

	if strings.TrimSpace(string(domain)) == "" {
		return nil, nil, errors.New("domain is empty")
	}

	routes, err := routerOp.Routes()
	if err != nil {
		// TODO: get routes later

		return nil, nil, errors.Wrap(err, onNew+"can't get routes")
	}

	handlers := map[[2]identity.Key]packs.Handler{}

	transpOp := transportHTTP{
		packsOp:  packsOp,
		routerOp: routerOp,

		routes: routes,
		domain: domain,
		path:   strconv.FormatInt(time.Now().UnixNano(), 10),

		handlers: handlers,
		mutex:    &sync.RWMutex{},
	}

	return &transpOp, transpOp.receiveEndpoint(), nil
}

const onSendOnly = "on transportHTTP.Send(): "

func (transpOp *transportHTTP) SendOnly(pack *packs.Pack, to identity.Key) (inPack *packs.Pack, doneAt *time.Time, err error) {
	if pack == nil {
		return nil, nil, errors.New(onSendOnly + "nothing to send")
	}

	target := to.Identity()
	if target == nil {
		return nil, nil, errors.Errorf(onSendOnly+"no target to send: '%s'", to)
	}

	var errs common.Errors

	route, ok := transpOp.routes[target.Domain]
	if !ok {
		errs = append(errs, errors.Errorf(onSendOnly+"no route to send: '%s'", to))
		return nil, nil, errs.Err()
	}

	url := route.URL()
	if url == "" {
		errs = append(errs, errors.Errorf(onSendOnly+"wrong route to send (empty .URL()): '%s'", to))
		return nil, nil, errs.Err()
	}

	packBytes, err := json.Marshal(pack)
	if err != nil {
		errs = append(errs, errors.Wrapf(err, onSendOnly+"can't marshal pack to send: '%#v'", pack))
		return nil, nil, errs.Err()
	}

	now := time.Now()
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(packBytes))
	if err != nil {
		errs = append(errs, errors.Wrapf(err, onSendOnly+"can't send to: '%s'", url))

		return nil, &now, errs.Err()
	}

	if resp.StatusCode != http.StatusOK {
		errs = append(errs, errors.Errorf(onSendOnly+"can't send to %s: status = %s", url, resp.Status))
		return nil, &now, errs.Err()
	}

	defer resp.Body.Close()

	inPack = &packs.Pack{} // TODO: is it necessarily???
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errs = append(errs, errors.Wrap(err, onSendOnly+"can't read inPack.Body"))
	}

	err = json.Unmarshal(bodyBytes, inPack)
	if err != nil {
		errs = append(errs, errors.Wrapf(err, onSendOnly+"can't unmarshal inPack(%s)", bodyBytes))
	}

	return inPack, &now, errs.Err()
}

// TODO!!! be careful because "pack.History = ..." isn't a thread safe action

const onSend = "on transportHTTP.Send(): "

func (transpOp *transportHTTP) Send(outPack *packs.Pack) (sentKey identity.Key, inPack *packs.Pack, err error) {
	if outPack == nil {
		return "", nil, errors.New(onSend + "nothing to send")
	}

	if strings.TrimSpace(string(outPack.Key)) == "" {
		transpOp.id++
		item := identity.Item{
			Domain: transpOp.domain,
			Path:   transpOp.path,
			ID:     strconv.FormatUint(transpOp.id, 10),
		}

		outPack.Key = item.Key()
	}

	ignoreProblems := outPack.Options.IsTrue("ignore_problens")

	var errs common.Errors

	id, err := transpOp.packsOp.Save(outPack, nil)
	if err != nil {
		if !ignoreProblems {
			return "", nil, errors.Wrap(err, onSend+"can't .Save()")
		}
		errs = append(errs, err)
	}

	var actionKey crud.ActionKey

	inPack, doneAtPtr, err := transpOp.SendOnly(outPack, outPack.To)
	if err != nil {
		errs = append(errs, err)
		actionKey = transport.DidntSendKey
	} else {
		actionKey = transport.SentKey

	}

	var doneAt time.Time
	if doneAtPtr != nil {
		doneAt = *doneAtPtr
	}

	action := crud.Action{
		// Identity: nil,
		Key:    actionKey,
		DoneAt: doneAt,
	}

	if id == "" {
		outPack.History = append(outPack.History, action)

	} else if historyHew, err := transpOp.packsOp.AddHistory(id, crud.History{action}, nil); err != nil {
		errs = append(errs, err)
		outPack.History = append(outPack.History, action)

	} else {
		outPack.History = historyHew

	}

	return "", inPack, errs.Err()
}

const onHistory = "on transportHTTP.History(): "

func (transpOp *transportHTTP) History(packKey identity.Key) (trace []crud.Action, err error) {
	return nil, common.ErrNotImplemented
}

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
