package flow

import "time"

type Origin struct {
	Source string
	Key    string
	Time   *time.Time
	Data   interface{}
}
