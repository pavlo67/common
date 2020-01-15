package connector

import "github.com/pavlo67/partes/connector"

func NewListener(receiverOp Operator) (*Listener, error) {
	listener := &Listener{
		receiverOp: receiverOp,
	}
	return listener, listener.Reset()
}

type Listener struct {
	receiverOp Operator
	envelopes  []connector.Envelope
	sessionID  uint64
}

func (l Listener) Reset() error {
	var err error
	l.envelopes, l.sessionID, err = l.receiverOp.Envelopes(0)

	return err
}

func (l Listener) ReadNext() ([]connector.Message, error) {
	envelopes, messages, err := MessagesNew(l.receiverOp, l.envelopes, l.sessionID)
	l.envelopes = append(l.envelopes, envelopes...)

	return messages, err
}
