package news_datastore_mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	"github.com/pavlo67/partes/crud"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/flow"
	"github.com/pavlo67/punctum/flow/datastore"
	"github.com/pavlo67/punctum/processor.old/news"
)

var _ news.Operator = &flowMySQL{}

type flowMySQL struct {
	dsOp datastore.Operator
}

const onNew = "on news_mysql.New()"

func New(dsOp datastore.Operator) (*flowMySQL, error) {
	if dsOp == nil {
		return nil, errors.New(onNew + ": no datastore.Operator")
	}

	return &flowMySQL{dsOp}, nil
}

const onAdd = "on flowMySQL.Save()"

func (flowOp flowMySQL) Save(item *news.Item) error {
	if item == nil {
		return errors.Wrap(basis.ErrNull, onAdd)
	}

	dataItem := &flow.Item{
		Source:      item.Source,
		Original:    item.Original,
		ContentKey:  item.ContentKey,
		ContentType: news.ContentType,
		Content:     item.Content,
		Status:      item.Status,
		History:     item.History,
	}

	return flowOp.dsOp.Save(dataItem)
}

func (flowOp flowMySQL) KeyExists(key string) (bool, error) {
	return flowOp.dsOp.KeyExists(news.ContentType, key)
}

func (flowOp flowMySQL) LastKey(options *content.ListOptions) (string, error) {
	return flowOp.dsOp.LastKey(news.ContentType, options)
}

func (flowOp flowMySQL) ReadList(options *content.ListOptions) ([]news.Item, uint64, error) {
	var ok bool
	var items []news.Item
	dataItems, allCnt, err := flowOp.dsOp.ReadList(options)

	for _, dataItem := range dataItems {
		item := news.Item{
			ID:         dataItem.ID,
			Source:     dataItem.Source,
			Original:   dataItem.Original,
			ContentKey: dataItem.ContentKey,
			Status:     dataItem.Status,
			History:    dataItem.History,
			StoredAt:   dataItem.StoredAt,
		}

		if item.Content, ok = dataItem.Content.(*news.Content); !ok {
			return nil, allCnt, errors.Errorf("wrong dataItem.Contentus (%#v), expected flow.Contentus", dataItem.Content)
		}

		items = append(items, item)
	}

	return items, allCnt, err
}

func (flowOp flowMySQL) Delete(options *content.ListOptions) (crud.Result, error) {
	return flowOp.dsOp.DeleteList(options)
}

func (flowOp flowMySQL) Close() error {
	return flowOp.dsOp.Close()
}
