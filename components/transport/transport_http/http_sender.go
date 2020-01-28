package transport_http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/identity"
	"github.com/pavlo67/workshop/components/packs"
	"github.com/pavlo67/workshop/components/transport"
	"github.com/pkg/errors"
)

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

	l.Infof("SENDER sent: %s", packBytes)

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

const onSend = "on transportHTTP.Send(): "

func (transpOp *transportHTTP) Send(outPack *packs.Pack) (sentKey identity.Key, targetTaskID common.ID, err error) {
	if outPack == nil {
		return "", "", errors.New(onSend + "nothing to send")
	}

	if strings.TrimSpace(string(outPack.Key)) == "" {
		transpOp.id++
		item := identity.Item{
			Domain: transpOp.domain,
			Path:   transpOp.path,
			ID:     common.ID(strconv.FormatUint(transpOp.id, 10)),
		}

		outPack.Key = item.Key()
	}

	ignoreProblems := outPack.Options.IsTrue("ignore_problens")

	var errs common.Errors

	id, err := transpOp.packsOp.Save(outPack, nil)
	if err != nil {
		if !ignoreProblems {
			return "", "", errors.Wrap(err, onSend+"can't .Save()")
		}
		errs = append(errs, err)
	}

	var actionKey crud.ActionKey

	_, doneAtPtr, err := transpOp.SendOnly(outPack, outPack.To)
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
		// ActorKey: nil,
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

	// TODO!!! targetTaskID
	return "", "", errs.Err()
}

const onHistory = "on transportHTTP.History(): "

func (transpOp *transportHTTP) History(packKey identity.Key) (trace []crud.Action, err error) {
	return nil, common.ErrNotImplemented
}
