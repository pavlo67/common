package scheduler

import (
	"time"
)

func Run(interval time.Duration, startImmediately bool, task Task) {
	now := time.Now()
	delta := time.Duration(now.UnixNano() % int64(interval))

	var timeScheduled time.Time

	if startImmediately {
		timeScheduled = now.Add(-delta)
	} else {
		timeScheduled = now.Add(interval - delta)
	}

	for {
		rest := timeScheduled.Sub(time.Now())
		if rest > 0 {
			l.Infof("%s: next scheduled task run.", timeScheduled.Format(time.RFC3339))
			time.Sleep(rest)
			continue
		}

		if rest > -interval {
			l.Infof("%s: task (%s) started...", timeScheduled.Format(time.RFC3339), task.Name())

			err := task.Run(timeScheduled)
			if err != nil {
				l.Errorf("on task(%s).Run(): %s", task.Name(), err)
			}

			l.Infof("%s: task (%s) finished", time.Now().Format(time.RFC3339), task.Name())
		}

		timeScheduled = timeScheduled.Add(interval)
	}
}
