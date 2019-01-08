package founts_leveldb

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/crud"
	"github.com/pavlo67/punctum/processor"
	"github.com/pavlo67/punctum/processor/founts"
)

type fountsLevelDB struct {
	path string
	db   *leveldb.DB
}

const onNew = "on founts_leveldb.New()"

func New(path string) (*fountsLevelDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "can't leveldb.OpenFile('%s', nil)", path)
	}

	return &fountsLevelDB{
		path: path,
		db:   db,
	}, nil
}

const onSave = "on fountsLevelDB.Save()"

func (fountsOp *fountsLevelDB) Save(url string, logItems ...processor.LogItem) error {
	var item founts.Item

	dataReaded, err := fountsOp.db.Get([]byte(url), nil)
	if err == leveldb.ErrNotFound {
		item.CreatedAtStr = time.Now().Format(time.RFC3339)
		// ok
	} else if err != nil {
		return errors.Wrapf(err, onSave+": can't fountsOp.db.Get('%s', nil)", url)
	} else {
		err = json.Unmarshal(dataReaded, &item)
		if err != nil {
			return errors.Wrapf(err, onSave+": can't json.Unmarshal('%s', &founts.Item)", dataReaded)
		}
		item.URL = ""
	}

	item.Log = append(item.Log, logItems...)

	dataToSave, err := json.Marshal(item)
	if err != nil {
		return errors.Wrapf(err, onSave+": can't json.Marshal(%#v)", dataToSave)
	}

	err = fountsOp.db.Put([]byte(url), dataToSave, nil)
	if err != nil {
		return errors.Wrapf(err, onSave+": can't fountsOp.db.Put('%s', '%s', nil)", url, dataToSave)
	}

	return nil
}

const onRead = "on fountsLevelDB.Read()"

func (fountsOp *fountsLevelDB) Read(url string) (*founts.Item, error) {
	dataReaded, err := fountsOp.db.Get([]byte(url), nil)

	if err == leveldb.ErrNotFound {
		return nil, err
	} else if err != nil {
		return nil, errors.Wrapf(err, onRead+": can't fountsOp.db.Get('%s', nil)", url)
	}

	var item founts.Item
	err = json.Unmarshal(dataReaded, &item)
	if err != nil {
		return nil, errors.Wrapf(err, onRead+": can't json.Unmarshal('%s', &founts.Item)", dataReaded)
	}

	item.URL = url

	// log.Printf("3: %#v", item)

	return &item, nil
}

const onReadList = "on fountsLevelDB.ReadList()"

func (fountsOp *fountsLevelDB) ReadList(*crud.ReadOptions) ([]founts.Item, *uint64, error) {
	var items []founts.Item
	var errs basis.Errors

	iter := fountsOp.db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		dataReaded := iter.Value()

		var item founts.Item
		err := json.Unmarshal(dataReaded, &item)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, onReadList+": can't json.Unmarshal('%s', &founts.Item)", dataReaded))
			continue
		}

		item.URL = string(key)
		items = append(items, item)
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		errs = append(errs, errors.Wrap(err, onReadList+": can't iter.Release()"))
	}

	return items, nil, errs.Err()
}

const onDelete = "on fountsLevelDB.DeleteList()"

func (fountsOp *fountsLevelDB) Delete(url string) error {
	err := fountsOp.db.Delete([]byte(url), nil)
	if err != nil {
		return errors.Wrapf(err, onDelete+": can't fountsOp.db.Delete('%s', nil)", url)
	}

	return nil
}

func (fountsOp *fountsLevelDB) Close() error {
	return errors.Wrap(fountsOp.db.Close(), "on fountsLevelDB.db.Close()")
}

const onClean = "on fountsLevelDB.clean()"

func (fountsOp *fountsLevelDB) clean() error {
	var errs basis.Errors

	iter := fountsOp.db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		err := fountsOp.db.Delete(key, nil)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, onClean+": can't fountsOp.db.Delete('%s', nil)", key))
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
