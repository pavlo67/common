package connector

import (
	"time"

	"github.com/pavlo67/punctum/things_old/files"
)

type Envelope struct {
	ID   string
	Size uint64
}

type Message struct {
	Envelope
	From, To, CC, Subject, ContentEncoding string
	Date                                   time.Time
	Body                                   string
	Files                                  []files.Item
}
