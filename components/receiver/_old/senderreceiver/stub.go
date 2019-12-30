package senderreceiver

import (
	"github.com/pavlo67/partes/connector"
	"github.com/pavlo67/partes/connector/receiver"
	"github.com/pavlo67/partes/connector/sender"
	"github.com/pavlo67/punctum/basis"
)

var _ sender.Operator = &senderreceiverstub{}

func New() (*senderreceiverstub, error) {
	return &senderreceiverstub{}, nil
}

type senderreceiverstub struct {
	messages []connector.Message
}

func (ss *senderreceiverstub) Send(message connector.Message) error {
	ss.messages = append(ss.messages, message)

	return nil
}

var _ receiver.Operator = &senderreceiverstub{}

func (ss senderreceiverstub) Start() (sessionID uint64, err error) {
	return 0, nil
}

func (ss senderreceiverstub) Stop(sessionID uint64) {

}

func (ss senderreceiverstub) URL() string {
	return ""
}

func (ss senderreceiverstub) Envelopes(sessionID0 uint64) (envelopesList []connector.Envelope, sessionID uint64, err error) {
	var envelopes []connector.Envelope

	for _, m := range ss.messages {
		envelopes = append(envelopes, m.Envelope)
	}

	return envelopes, 0, nil
}

func (ss senderreceiverstub) Read(id string, sessionID0 uint64) (message *connector.Message, sessionID uint64, err error) {
	for _, m := range ss.messages {
		if m.ID == id {
			return &m, 0, nil
		}

		ss.messages = append(ss.messages, m)
	}

	return nil, 0, basis.ErrNotFound
}

func (ss senderreceiverstub) ReadRaw(id string, sessionID0 uint64) (*string, uint64, error) {
	return nil, 0, basis.ErrNotImplemented
}

func (ss senderreceiverstub) Delete(id string, sessionID0 uint64) (sessionID uint64, err error) {
	messages := ss.messages

	ss.messages = nil
	for _, m := range messages {
		if m.ID == id {
			continue
		}

		ss.messages = append(ss.messages, m)
	}

	return 0, nil
}

func (ss senderreceiverstub) Close() error {
	return nil
}
