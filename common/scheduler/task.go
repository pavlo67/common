package scheduler

import "time"

type Task interface {
	Name() string
	Run(timeSheduled time.Time) error
}
