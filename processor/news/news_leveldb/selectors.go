package news_leveldb

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb/util"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/strlib"
	"github.com/pavlo67/punctum/crud"
)

var reDigits = regexp.MustCompile(`^\d+$`)

type CheckIter func(key, value []byte, exemplar interface{}) bool

const onConditions = "on seveldblib.RangesAndCheck()"

func RangesAndCheck(opt *content.ListOptions) (*util.Range, CheckIter, error) {
	if opt == nil || opt.Selector == nil {
		return nil, nil, nil
	}

	_, urls, timeMin, timeMax, err := conditions(*opt.Selector)
	if err != nil {
		return nil, nil, errors.Wrapf(err, onConditions+" conditionsUnary() with .Selector.First: %#v", opt.Selector.First)
	}

	var ranges *util.Range
	if len(urls) == 1 {
		url0 := strings.TrimSpace(urls[0])
		pos := strings.Index(url0, "#")
		if pos >= 0 {
			url0 = url0[:pos]
		}

		ranges = &util.Range{
			Start: []byte(url0 + "#"),
			Limit: []byte(url0 + "#9"),
		}
		urls = nil
	}

	if len(urls) < 1 && timeMin == nil && timeMax == nil {
		return ranges, nil, nil
	}

	var keyTimeMin, keyTimeMax uint64

	if timeMin != nil {
		keyTimeMin = uint64(timeMin.Unix())
	}
	if timeMax != nil {
		keyTimeMax = uint64(timeMax.Unix())
	}

	return ranges, func(key, value []byte, exemplar interface{}) bool {
		var keyURL string
		var keyTime uint64

		keyParts := strings.Split(string(key), "#")
		if len(keyParts) >= 1 {
			keyURL = keyParts[0]
			if len(keyParts) >= 2 {
				keyTime, _ = strconv.ParseUint(keyParts[1], 10, 64)
			}
		}

		return (len(urls) < 1 || strlib.In(urls, keyURL)) && keyTime >= keyTimeMin && (timeMax == nil || keyTime < keyTimeMax)
	}, nil

}

func conditions(term basis.Term) (kv []byte, urls []string, timeMin, timeMax *time.Time, err error) {
	kv, urls, timeMin, timeMax, err = conditionsUnary(term.First)
	if err != nil {
		return nil, nil, nil, nil, errors.Wrapf(err, "conditionsUnary() with .Selector.First: %#v", term.First)
	}

	for i, next := range term.Next {
		switch next.Operation {
		case basis.GE, basis.LT:
			if i > 0 {
				return nil, nil, nil, nil, errors.Wrapf(err, "too many operands to compare: %#v", term.Next)
			}

			kvNext, urlsNext, timeMinNext, timeMaxNext, err := conditionsUnary(next.TermUnary)
			if err != nil {
				return nil, nil, nil, nil, errors.Wrapf(err, "conditionsUnary() with .Next[%d]: %#v", i, next.TermUnary)
			}
			if urlsNext != nil || timeMinNext != nil || timeMaxNext != nil {
				return nil, nil, nil, nil, errors.Wrapf(err, "too complicate operands to compare: %#v", term.Next)
			}
			if crud.FieldKey(kv) != crud.TimeField {
				return nil, nil, nil, nil, errors.Wrapf(err, "wrong key to compare: %#v", term)
			}
			if string(kvNext) == "" {
				return nil, urls, timeMin, timeMax, nil
			}
			timeNext, err := time.Parse(time.RFC3339, string(kvNext))
			if err != nil {
				return nil, nil, nil, nil, errors.Wrapf(err, "wrong time value to compare: %#v", term)
			}

			if next.Operation == basis.GE {
				return nil, urls, &timeNext, nil, nil
			}
			return nil, urls, nil, &timeNext, nil

		case basis.AND:
			_, urlsNext, timeMinNext, timeMaxNext, err := conditionsUnary(next.TermUnary)
			if err != nil {
				return nil, nil, nil, nil, errors.Wrapf(err, "conditionsUnary() with .Next[%d]: %#v", i, next.TermUnary)
			}
			if timeMinNext != nil {
				if timeMin != nil {
					return nil, nil, nil, nil, errors.Wrapf(err, "duplicate timeMin in selector: %#v", term)
				}
				timeMin = timeMinNext
			}
			if timeMaxNext != nil {
				if timeMax != nil {
					return nil, nil, nil, nil, errors.Wrapf(err, "duplicate timeMax in selector: %#v", term)
				}
				timeMax = timeMaxNext
			}
			urls = append(urls, urlsNext...)

		default:
			return nil, nil, nil, nil, errors.Errorf("forbidden operation: %s", next.Operation)
		}
	}

	return nil, urls, timeMin, timeMax, nil
}

func conditionsUnary(termUnary basis.TermUnary) (kv []byte, urls []string, timeMin, timeMax *time.Time, err error) {
	if termUnary.OperationUnary != nil {
		return nil, nil, nil, nil, errors.New(".OperationUnary doesn't implemented yet")
	}

	switch val := termUnary.Value.(type) {
	case string:
		return []byte(val), nil, nil, nil, nil
	case basis.Term:
		return conditions(val)
	case basis.TermUnary:
		return conditionsUnary(val)
	case basis.TermOneOfStr:
		if crud.FieldKey(val.Key) != crud.URLField {
			return nil, nil, nil, nil, errors.Errorf("wrong key (%s) for basis.TermOneOf", val.Key)
		}
		return nil, val.Values, nil, nil, nil
	}

	return nil, nil, nil, nil, errors.Errorf("TermUnary.Value has a wrong type: %#v", termUnary.Value)
}
