package transport_smtp

import (
	"crypto/tls"

	"strings"

	"github.com/pavlo67/partes/connector"
	"github.com/pavlo67/partes/connector/sender"
	"github.com/pavlo67/punctum/starter/config"
)

func New(smtpConfig config.ServerAccess, senderConfig map[string]string) (sender.Operator, error) {
	return &senderSMTPGomail{
		smtpConfig:   smtpConfig,
		senderConfig: senderConfig,
	}, nil
}

var _ sender.Operator = &senderSMTPGomail{}

type senderSMTPGomail struct {
	smtpConfig   config.ServerAccess
	senderConfig map[string]string
}

func (sm *senderSMTPGomail) Send(message connector.Message) error {
	m := gomail.NewMessage()

	from := strings.TrimSpace(message.From)
	if from == "" {
		from = sm.senderConfig["from"]
	}
	m.SetHeader("From", from)

	m.SetHeader("To", message.To)
	if message.CC != "" {
		m.SetAddressHeader("Cc", message.CC, "")
	}

	// date :=

	m.SetHeader("Subject", message.Subject)
	m.SetBody("text/html", string(message.Body))

	for _, f := range message.Files {
		m.Attach(f.Name)
	}

	d := gomail.NewDialer(sm.smtpConfig.Host, sm.smtpConfig.Port, sm.smtpConfig.User, sm.smtpConfig.Pass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
