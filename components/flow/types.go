package flow

import "time"

type Origin struct {
	Source string     `bson:",omitempty"    json:",omitempty"`
	Key    string     `bson:",omitempty"    json:",omitempty"`
	Time   *time.Time `bson:",omitempty"    json:",omitempty"`
	Data   string     `bson:",omitempty"    json:",omitempty"`
}
