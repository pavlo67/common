package processor

import "time"

type LogItem struct {
	Started  time.Time
	Finished time.Time
	Success  bool
	Info     string
}
