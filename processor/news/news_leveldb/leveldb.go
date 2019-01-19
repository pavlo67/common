package news_leveldb

import (
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"

	"encoding/json"
	"time"

	"strconv"

	"strings"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/crud"
	"github.com/pavlo67/punctum/processor/flow"
	"github.com/pavlo67/punctum/processor/news"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type newsLevelDB struct {
	path string
	db   *leveldb.DB
}

const onNew = "on news_leveldb.New()"

func New(path string) (*newsLevelDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "can't leveldb.OpenFile('%s', nil)", path)
	}

	return &newsLevelDB{
		path: path,
		db:   db,
	}, nil
}

const onSave = "on newsLevelDB.Save()"

func (newsOp *newsLevelDB) Save(item *news.Item) error {
	if item == nil {
		return errors.Wrap(basis.ErrNull, onSave)
	}

	if item.SavedAt == nil {
		savedAt := time.Now()
		item.SavedAt = &savedAt
	}

	key := item.Key(strconv.FormatInt(item.SavedAt.Unix(), 10))
	if key == "" {
		return errors.Errorf(onSave+": empty key for item (%#v)", item)
	}

	dataToSave, err := json.Marshal(item)
	if err != nil {
		return errors.Wrapf(err, onSave+": can't json.Marshal(%#v)", dataToSave)
	}

	err = newsOp.db.Put([]byte(key), dataToSave, nil)
	if err != nil {
		return errors.Wrapf(err, onSave+": can't newsOp.db.Put('%s', '%s', nil)", key, dataToSave)
	}

	return nil
}

const onHas = "on newsLevelDB.Has()"

func (newsOp *newsLevelDB) Has(src *flow.Source) (bool, error) {
	keyStart := src.Key("0")
	if keyStart == "" {
		return false, errors.Errorf(onHas+": empty keyStart to check from source: %#v", src)
	}
	keyLimit := src.Key("9")
	url := strings.TrimSpace(src.URL)
	sourceID := strings.TrimSpace(src.SourceID)

	iter := newsOp.db.NewIterator(&util.Range{Start: []byte(keyStart), Limit: []byte(keyLimit)}, nil)
	for iter.Next() {
		keyParts := strings.Split(string(iter.Key()), "#")
		if len(keyParts) == 3 && keyParts[0] == url && keyParts[2] == sourceID {
			iter.Release()
			return true, nil
		}
	}
	iter.Release()

	return false, errors.Wrap(iter.Error(), onHas)
}

type CheckKey func([]byte) bool

func ranges(opt *crud.ReadOptions) (*util.Range, CheckKey) {

	if opt == nil || (opt.RangeMin == "" && opt.RangeMax == "") {
		return nil, nil
	}

	if opt.RangedBy == crud.TimeField {
		var rangeMin, rangeMax uint64
		var checkMax bool

		if opt.RangeMin != "" {
			rangeMin, _ = strconv.ParseUint(opt.RangeMin, 10, 64)
		}
		if opt.RangeMax != "" {
			checkMax = true
			rangeMax, _ = strconv.ParseUint(opt.RangeMax, 10, 64)
		}

		return nil, func(key []byte) bool {
			keyParts := strings.Split(string(key), "#")
			var keyTime uint64
			if len(keyParts) >= 2 {
				keyTime, _ = strconv.ParseUint(keyParts[1], 10, 64)
			}

			return keyTime >= rangeMin && (!checkMax || keyTime < rangeMax)
		}
	}

	var ranges *util.Range
	var checkKey CheckKey
	if opt.RangeMin != "" {
		if opt.RangeMax != "" {
			ranges = &util.Range{Start: []byte(opt.RangeMin), Limit: []byte(opt.RangeMax)}
		} else {
			ranges = &util.Range{Start: []byte(opt.RangeMin)}
		}
	} else {
		// opt.RangeMax != "" due to: if opt == nil || (opt.RangeMin == "" && opt.RangeMax == "")
		ranges = &util.Range{Limit: []byte(opt.RangeMax)}
	}

	return ranges, checkKey
}

const onReadList = "on newsLevelDB.ReadList()"

func (newsOp *newsLevelDB) ReadList(opt *crud.ReadOptions) ([]news.Item, *uint64, error) {

	var items []news.Item
	var errs basis.Errors

	ranges, checkKey := ranges(opt)

	iter := newsOp.db.NewIterator(ranges, nil)
	for iter.Next() {
		if checkKey != nil && !checkKey(iter.Key()) {
			continue
		}

		value := iter.Value()

		var item news.Item
		err := json.Unmarshal(value, &item)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, onReadList+": can't json.Unmarshal('%s', &news.Item)", value))
			continue
		}
		items = append(items, item)
	}
	iter.Release()

	return items, nil, errors.Wrap(iter.Error(), onReadList)
}

const onDeleteList = "on newsLevelDB.DeleteList()"

func (newsOp *newsLevelDB) DeleteList(opt *crud.ReadOptions) error {
	var errs basis.Errors

	ranges, checkKey := ranges(opt)

	iter := newsOp.db.NewIterator(ranges, nil)
	for iter.Next() {
		key := iter.Key()

		if checkKey != nil && !checkKey(key) {
			continue
		}

		err := newsOp.db.Delete(key, nil)
		if err != nil {
			errs = append(errs, err)
		}

	}
	iter.Release()

	return errors.Wrap(errs.Append(iter.Error()).Err(), onDeleteList)
}

func (newsOp *newsLevelDB) Close() error {
	return errors.Wrap(newsOp.db.Close(), "on newsLevelDB.db.Close()")
}

const onClean = "on newsLevelDB.clean()"

func (newsOp *newsLevelDB) clean() error {
	var errs basis.Errors

	iter := newsOp.db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		err := newsOp.db.Delete(key, nil)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, onClean+": can't newsOp.db.Delete('%s', nil)", key))
			continue
		}
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		errs = append(errs, errors.Wrap(err, onClean+": can't iter.Release()"))
	}

	return errs.Err()
}
