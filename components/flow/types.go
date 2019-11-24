package flow

import "time"

type Origin struct {
	Source     string
	Key        string
	OriginTime time.Time
	OriginData []byte
}
