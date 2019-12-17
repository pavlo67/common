package tagger_sqlite

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/libraries/sqllib"
	"github.com/pavlo67/workshop/common/libraries/strlib"
)

const onCountTagsOnTag = "on taggerSQLite.countTag(): "

func (taggerOp *taggerSQLite) countTag(tagLabel string, passedTags []string, labelsRemoved []string, stmCountTags, stmListTags, stmAddTag *sql.Stmt) ([]string, error) {
	if strlib.In(passedTags, tagLabel) {
		return passedTags, nil
	}

	var partedSize uint64
	partedSizePtr := &partedSize
	values := []interface{}{taggerOp.ownInterfaceKey, tagLabel}
	row := stmCountTags.QueryRow(values...)
	if err := row.Scan(&partedSizePtr); err != nil && err != sql.ErrNoRows {
		return passedTags, errors.Wrapf(err, onCountTagsOnTag+": can't tx.QueryRow(%s, %#v)", taggerOp.sqlCountTag, values)
	}

	var labelsOnTag []string
	values = []interface{}{taggerOp.ownInterfaceKey, tagLabel}
	rows, err := stmListTags.Query(values...)
	if err != sql.ErrNoRows && err != nil {
		return passedTags, errors.Wrapf(err, onCountTagsOnTag+": can't tx.Query(%s, %#v)", taggerOp.sqlListTags, values)
	}
	defer rows.Close()
	for rows.Next() {
		var tagLabel, relation string
		err = rows.Scan(&tagLabel, &relation)
		if err != nil {
			return passedTags, errors.Wrapf(err, onCountTagsOnTag+": can't tx.ScanQueryRow(%s, %#v)", taggerOp.sqlListTags, values)
		}
		labelsOnTag = append(labelsOnTag, tagLabel)
	}
	err = rows.Err()
	if err != nil {
		return passedTags, errors.Wrapf(err, onCountTagsOnTag+": "+sqllib.RowsError, taggerOp.sqlListTags, values)
	}

	values = []interface{}{tagLabel, len(labelsOnTag) > 0, partedSize}
	if _, err := stmAddTag.Exec(values...); err != nil {
		return passedTags, errors.Wrapf(err, onCountTagsOnTag+": can't tx.Exec(%s, %#v)", taggerOp.sqlAddTag, values)
	}
	// TODO: don't forget! this must be done before loop with (top!) labelsToCount
	passedTags = append(passedTags, tagLabel)

	labelsToCount := labelsOnTag
	for _, labelRemoved := range labelsRemoved {
		if !strlib.In(labelsToCount, labelRemoved) {
			labelsToCount = append(labelsToCount, labelRemoved)
		}
	}

	for _, labelToCount := range labelsToCount {
		if passedTags, err = taggerOp.countTag(labelToCount, passedTags, nil, stmCountTags, stmListTags, stmAddTag); err != nil {
			return passedTags, errors.Wrapf(err, "on tag '%s'", tagLabel)
		}
	}

	return passedTags, nil
}

const onCountTagsChanged = "on taggerSQLite.countTagChanged(): "

func (taggerOp *taggerSQLite) countTagChanged(key joiner.InterfaceKey, id common.ID, tagLabelsRemoved []string, tx *sql.Tx) error {
	if key != taggerOp.ownInterfaceKey {
		return nil
	}
	tagLabel := string(id)

	stmCountTags, err := tx.Prepare(taggerOp.sqlCountTag)
	if err != nil {
		return errors.Wrapf(err, onCountTagsChanged+": can't tx.Prepare(%s)", taggerOp.sqlCountTag)
	}
	stmListTags, err := tx.Prepare(taggerOp.sqlListTags)
	if err != nil {
		return errors.Wrapf(err, onCountTagsChanged+": can't tx.Prepare(%s)", taggerOp.sqlListTags)
	}
	stmAddTag, err := tx.Prepare(taggerOp.sqlAddTag)
	if err != nil {
		return errors.Wrapf(err, onCountTagsChanged+": can't tx.Prepare(%s)", taggerOp.sqlAddTag)
	}

	if _, err := taggerOp.countTag(tagLabel, nil, tagLabelsRemoved, stmCountTags, stmListTags, stmAddTag); err != nil {
		return errors.Wrapf(err, onCountTagsChanged+": can't taggerOp.countTag(%s, ...)", tagLabel)
	}

	return nil
}
