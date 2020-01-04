package sender_http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/identity"

	"io/ioutil"

	"github.com/pavlo67/workshop/components/packs"
	"github.com/pavlo67/workshop/components/router"
	"github.com/pavlo67/workshop/components/sender"
)

var _ sender.Operator = &senderHTTP{}

type senderHTTP struct {
	packsOp  packs.Operator
	routerOp router.Operator

	routes router.Routes
}

const onNew = "on sender_http.New(): "

func New(packsOp packs.Operator, routerOp router.Operator) (sender.Operator, error) {
	if packsOp == nil {
		return nil, errors.New(onNew + "no packs.Operator")
	}

	if routerOp == nil {
		return nil, errors.New(onNew + "no router.Operator")
	}

	routes, err := routerOp.Routes()
	if err != nil {
		// TODO: get routes later

		return nil, errors.Wrap(err, onNew+"can't get routes")
	}

	senderOp := senderHTTP{
		packsOp:  packsOp,
		routerOp: routerOp,
		routes:   routes,
	}

	return &senderOp, nil
}

const onSendOnly = "on senderHTTP.SendOne(): "

func (senderOp *senderHTTP) SendOnly(pack *packs.Pack, to identity.Key) (inPack *packs.Pack, doneAt *time.Time, err error) {
	if pack == nil {
		return nil, nil, errors.New(onSendOnly + "nothing to send")
	}

	target := to.Identity()
	if target == nil {
		return nil, nil, errors.Errorf(onSendOnly+"no target to send: '%s'", to)
	}

	var errs common.Errors

	route, ok := senderOp.routes[target.Domain]
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

const onSendOne = "on senderHTTP.SendOne(): "

func (senderOp *senderHTTP) SendOne(outPack *packs.Pack, to identity.Key, ignoreProblems bool) (inPack *packs.Pack, err error) {
	if outPack == nil {
		return nil, errors.New(onSendOne + "nothing to send")
	}

	var errs common.Errors

	id, err := senderOp.packsOp.Save(outPack, nil)
	if err != nil {
		if !ignoreProblems {
			return nil, errors.Wrap(err, onSendOne+"can't .Save()")
		}
		errs = append(errs, err)
	}

	var actionKey crud.ActionKey

	inPack, doneAtPtr, err := senderOp.SendOnly(outPack, to)
	if err != nil {
		errs = append(errs, err)
		actionKey = sender.DidntSendKey
	} else {
		actionKey = sender.SentKey

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

	} else if historyHew, err := senderOp.packsOp.AddHistory(id, crud.History{action}, nil); err != nil {
		errs = append(errs, err)
		outPack.History = append(outPack.History, action)

	} else {
		outPack.History = historyHew

	}

	return inPack, errs.Err()
}

const onSend = "on senderHTTP.Send(): "

func (senderOp *senderHTTP) Handle(pack *packs.Pack) (*packs.Pack, error) {
	//if pack == nil {
	//	return nil, errors.New(onSend + "nothing to send")
	//}
	//
	//var errs common.Errors
	//
	//for _, to := range pack.To {
	//	senderOp.SendOne(Pack, to)
	//}

	return nil, common.ErrNotImplemented
}

const onHistory = "on senderHTTP.History(): "

func (senderOp *senderHTTP) History(packKey identity.Key) (trace []crud.Action, err error) {
	return nil, common.ErrNotImplemented
}
