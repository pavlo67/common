package receiverpop3

import (
	"io/ioutil"
	"log"
	"net/mail"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/pavlo67/partes/connector"
	"github.com/pavlo67/partes/connector/receiver"
	"github.com/pavlo67/partes/validator"
	"github.com/pavlo67/punctum/starter/config"
)

func NewPOP3BytBox(serverAccessConfig config.ServerAccess, emailAddress string) (receiver.Operator, error) {
	return &receiverPOP3BytBox{
		serverAccessConfig: serverAccessConfig,
		emailAddress:       emailAddress,
		mutex:              &sync.Mutex{},
		clients:            map[uint64]*pop3.Client{},
	}, nil
}

type receiverPOP3BytBox struct {
	serverAccessConfig config.ServerAccess
	emailAddress       string
	mutex              *sync.Mutex
	clients            map[uint64]*pop3.Client
	sessionID          uint64
}

var _ receiver.Operator = &receiverPOP3BytBox{}

var noOp = errors.New("receiver.Operator == nil	")

func (recOp *receiverPOP3BytBox) start() (*pop3.Client, uint64, error) {
	if recOp == nil {
		return nil, 0, errors.Wrap(noOp, "on recOp.start()")
	}

	client, err := pop3.DialTLS(recOp.serverAccessConfig.Host)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "on pop3.DialTLS(%s)", recOp.serverAccessConfig.Host)
	}

	err = client.Auth(recOp.serverAccessConfig.User, recOp.serverAccessConfig.Pass)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "on recOp.client.Auth(%s, %s)", recOp.serverAccessConfig.User, recOp.serverAccessConfig.Pass)
	}

	recOp.mutex.Lock()
	defer recOp.mutex.Unlock()

	recOp.sessionID++
	recOp.clients[recOp.sessionID] = client

	return client, recOp.sessionID, nil
}

func (recOp *receiverPOP3BytBox) client(sessionID uint64) (*pop3.Client, uint64, error) {
	if sessionID != 0 {
		recOp.mutex.Lock()
		defer recOp.mutex.Unlock()

		if client, ok := recOp.clients[sessionID]; ok {
			return client, sessionID, nil
		}
		log.Print("on recOp.client(%d): no client for sessionID", sessionID)
	}

	return recOp.start()
}

func (recOp *receiverPOP3BytBox) Start() (uint64, error) {
	if recOp == nil {
		return 0, errors.Wrap(noOp, "on recOp.Start()")
	}

	_, sessionID, err := recOp.start()

	return sessionID, err
}

func (recOp *receiverPOP3BytBox) Stop(sessionID uint64) {
	if recOp == nil {
		log.Print("on recOp.Stop(%d): $s", sessionID, noOp)
		return
	}

	recOp.mutex.Lock()
	defer recOp.mutex.Unlock()

	if client, ok := recOp.clients[sessionID]; ok {
		delete(recOp.clients, sessionID)

		err := client.Quit()
		if err != nil {
			log.Print("on recOp.client.Quit(%d): $s", sessionID, err)
		}

		return
	}

	log.Printf("recOp.client.Quit(%d): no client for id", sessionID)
}

func (recOp *receiverPOP3BytBox) Close() error {
	if recOp != nil {
		for sessionID := range recOp.clients {
			recOp.Stop(sessionID)
		}
	}

	return nil
}

func (recOp *receiverPOP3BytBox) URL() string {
	return "mailto:" + recOp.emailAddress
}

func (recOp *receiverPOP3BytBox) Envelopes(sessionID0 uint64) ([]connector.Envelope, uint64, error) {
	client, sessionID, err := recOp.client(sessionID0)
	if err != nil {
		return nil, 0, err
	}

	var envelopes []connector.Envelope

	msgs, sizes, err := client.ListAll()
	for i, msg := range msgs {
		if i >= len(sizes) {
			envelopes = append(envelopes, connector.Envelope{ID: strconv.Itoa(msg)})
		} else {
			envelopes = append(envelopes, connector.Envelope{ID: strconv.Itoa(msg), Size: uint64(sizes[i])})
		}
	}

	return envelopes, sessionID, nil
}

func (recOp *receiverPOP3BytBox) ReadRaw(id string, sessionID0 uint64) (*string, uint64, error) {
	client, sessionID, err := recOp.client(sessionID0)
	if err != nil {
		return nil, 0, err
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, 0, errors.Errorf("on .ReadRaw: wrong id: '%s'", id)
	}

	msg, err := client.Retr(idInt)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "on client.Retr(%d)", idInt)
	}

	return &msg, sessionID, nil
}

func (recOp *receiverPOP3BytBox) Read(id string, sessionID0 uint64) (*connector.Message, uint64, error) {
	msg, sessionID, err := recOp.ReadRaw(id, sessionID0)
	if err != nil {
		return nil, 0, err
	}

	r := strings.NewReader(*msg)
	m, err := mail.ReadMessage(r)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "on strings.NewReader(&s)", *msg)
	}

	header := m.Header

	message := &connector.Message{
		Envelope:        connector.Envelope{ID: id, Size: uint64(len(*msg))},
		From:            header.Get("From"),
		To:              header.Get("To"),
		CC:              header.Get("CC"),
		Date:            validator.DateFromString(header.Get("Date"), "").Time(),
		Subject:         header.Get("Subject"),
		ContentEncoding: header.Get("Contentus-Encoding"),
	}

	// TODO: read files

	body, err := ioutil.ReadAll(m.Body)
	if err != nil {
		return message, sessionID, errors.Wrapf(err, "on ioutil.ReadAll(m.Body) from: %s", *msg)
	}
	message.Body = string(body)

	return message, sessionID, nil
}

func (recOp *receiverPOP3BytBox) Delete(id string, sessionID0 uint64) (uint64, error) {
	client, sessionID, err := recOp.client(sessionID0)
	if err != nil {
		return 0, err
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return 0, errors.Errorf("on .DeleteList: wrong id: '%s'", id)
	}

	err = client.Dele(idInt)
	if err != nil {
		return 0, errors.Wrapf(err, "om client.Dele(%d)", idInt)
	}

	return sessionID, nil
}
