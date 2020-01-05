package connector

import (
	"github.com/pavlo67/partes/connector"
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/starter/joiner"
	"github.com/pkg/errors"
)

const InterfaceKey joiner.InterfaceKey = "receiver"

type Operator interface {
	Start() (sessionID uint64, err error)
	Stop(sessionID uint64)
	URL() string

	Envelopes(sessionID0 uint64) (envelopesList []connector.Envelope, sessionID uint64, err error)
	Read(id string, sessionID0 uint64) (message *connector.Message, sessionID uint64, err error)

	ReadRaw(id string, sessionID0 uint64) (*string, uint64, error)
	Delete(id string, sessionID0 uint64) (sessionID uint64, err error)

	Close() error
}

func EnvelopesNew(recOp Operator, eOld []connector.Envelope, sessionID uint64) ([]connector.Envelope, error) {
	eCurrent, _, err := recOp.Envelopes(sessionID)
	if err != nil {
		return nil, errors.Wrap(err, "on .EnvelopesList()")
	}

	var eNew []connector.Envelope

E_CURRENT:
	for _, e := range eCurrent {
		for _, e0 := range eOld {
			if e.ID == e0.ID {
				continue E_CURRENT
			}
		}

		eNew = append(eNew, e)
	}

	return eNew, nil
}

func MessagesNew(recOp Operator, eOld []connector.Envelope, sessionID uint64) ([]connector.Envelope, []connector.Message, error) {
	envelopes, err := EnvelopesNew(recOp, eOld, sessionID)
	if err != nil {
		return nil, nil, errors.Wrap(err, "on EnvelopesListNew()")
	}

	var messages []connector.Message
	var errs basis.Errors

	for _, e := range envelopes {
		message, _, err := recOp.Read(e.ID, sessionID)
		if err != nil {
			errs = errs.Append(errors.Wrapf(err, "on .ReceiveByID(%s)", e.ID))
		} else if message == nil {
			errs = errs.Append(errors.Errorf("empty message on .ReceiveByID(%s)", e.ID))
		} else {
			messages = append(messages, *message)
		}
	}

	return envelopes, messages, errs.Err()
}
