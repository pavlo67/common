package config

import (
	"fmt"

	"github.com/pkg/errors"
)

type Listener struct {
	ID     string   `json:"id,omitempty"`
	Events []string `json:"events,omitempty"`
	Action string   `json:"action,omitempty"`
}

func readListener(key string, l0 interface{}) (*Listener, error) {
	if l, ok := l0.(string); ok {
		return &Listener{ID: l, Events: []string{"click"}, Action: key}, nil
	}

	var l1 []interface{}
	var ok bool

	if l1, ok = l0.([]interface{}); ok {
	} else if l00, ok := l0.([]string); ok {
		for _, l := range l00 {
			l1 = append(l1, l)
		}
	} else {
		l1 = append(l1, l0)
	}

	lst := new(Listener)
	lst.Action = key

	if len(l1) < 2 {
		lst.Events = []string{"click"}
	} else if event, ok := l1[1].(string); ok {
		lst.Events = []string{event}
	} else if events, ok := l1[1].([]string); ok {
		lst.Events = append(lst.Events, events...)
	} else if eventsTmp, ok := l1[1].([]interface{}); ok {
		events, err := stringifySlice(eventsTmp)
		if err != nil {
			return nil, errors.Wrapf(err, "bad listener JSON: %#v", l0)
		}
		lst.Events = append(lst.Events, events...)
	} else {
		return nil, fmt.Errorf("bad listener JSON: %#v", l0)
	}

	if len(l1) >= 1 {
		lst.ID, _ = l1[0].(string)
	}
	if lst.ID == "" {
		lst.ID = key
	}

	return lst, nil
}
