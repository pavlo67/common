package data_mongodb

import (
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/libraries/mgolib"
	"github.com/pavlo67/workshop/common/selectors"
	"github.com/pavlo67/workshop/components/data"
)

// TODO: use common context.Session!!!

// storage -----------------------------------------------------------------------------------------------------------------

var _ data.Operator = &dataMongoDB{}

type dataMongoDB struct {
	client     *mongo.Client
	collection *mongo.Collection
}

const onNewData = "on crud_mongodb.NewData()"

func NewData(access *config.Access, timeout time.Duration, dbName, collectionName string, exemplar data.Item) (data.Operator, crud.Cleaner, *mongo.Client, error) {
	if exemplar.Details == nil {
		return nil, nil, nil, errors.New("no exemplar.Details")
	}

	client, err := mgolib.Connect(access, timeout)
	if err != nil {
		return nil, nil, nil, err
	}
	database := client.Database(dbName)
	mgoOp := &dataMongoDB{
		client:     client,
		collection: database.Collection(collectionName),
	}

	return mgoOp, mgoOp, client, nil
}

// operator ----------------------------------------------------------------------------------------------------------------

//const onExemplar = "on dataMongoDB.Exemplar()"
//
//func (mgoOp dataMongoDB) Exemplar() crud.Item {
//	if reflect.TypeOf(mgoOp.exemplar.Details).Kind() == reflect.Ptr {
//		// Pointer:
//		return crud.Item{Details: reflect.New(reflect.ValueOf(mgoOp.exemplar.Details).Elem().Type()).Interface()}
//	}
//
//	// Not pointer:
//	return crud.Item{Details: reflect.New(reflect.TypeOf(mgoOp.exemplar.Details)).Elem().Interface()}
//}

const onSave = "on dataMongoDB.Save()"

func (mgoOp dataMongoDB) Save(items []data.Item, options *crud.SaveOptions) ([]common.ID, error) {
	// TODO: use Upsert

	var err error
	var ids []common.ID

	for i, item := range items {
		if item.ID != "" {
			id, err := primitive.ObjectIDFromHex(string(item.ID))
			if err != nil {
				return nil, errors.Wrapf(err, onSave+": can't primitive.ObjectIDFromHex(string(%s))", item.ID)
			}

			filter := bson.M{"_id": id}
			_, err = mgoOp.collection.DeleteMany(nil, filter)
			if err != nil {
				return nil, errors.Wrapf(err, onSave+": can't .DeleteMany(nil, %#v, nil)", filter)
			}

		}

		items[i].DetailsRaw, err = bson.Marshal(item.Details)
		if err != nil {
			return nil, errors.Wrapf(err, onSave+": can't bson.Marshal(item.Details = %#v)", item.Details)
		}

		res, err := mgoOp.collection.InsertOne(nil, &item)
		if err != nil {
			return nil, errors.Wrapf(err, onSave+": can't .InsertOne(nil, %#v)", item)
		}

		var id common.ID
		if objectID, ok := res.InsertedID.(primitive.ObjectID); ok {
			id = common.ID(objectID.Hex())

			// TODO!!! save tags

		} else {
			return ids, errors.Errorf("??? res.InsertedID = %#v", res.InsertedID)
		}

		ids = append(ids, id)

	}

	return ids, nil
}

const onRead = "on dataMongoDB.Read()"

func (mgoOp dataMongoDB) Read(id common.ID, options *crud.GetOptions) (*data.Item, error) {
	objectID, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return nil, errors.Wrapf(err, onSave+": can't primitive.ObjectIDFromHex(string(%s))", id)
	}

	filter := bson.M{"_id": objectID}
	res := mgoOp.collection.FindOne(nil, filter)
	if res.Err() != nil {
		if res.Err().Error() == mongo.ErrNoDocuments.Error() {
			return nil, nil
		}

		return nil, errors.Wrapf(res.Err(), onRead+": no result is returned with collection.FindOne(nil, %#v)", filter)
	}
	if res == nil {
		return nil, nil
	}

	item := data.Item{}

	err = res.Decode(&item)
	if err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			return nil, nil
		}

		return nil, errors.Wrapf(err, onRead+": on res.Decode(%#v)", res)
	}

	return &item, nil
}

const onDetails = "on dataMongoDB.Details()"

func (mgoOp dataMongoDB) Details(item *data.Item, exemplar interface{}) error {
	if item == nil {
		return errors.Wrap(common.ErrNull, onDetails+": item == nil")
	}
	if len(item.DetailsRaw) < 1 {
		return errors.Wrap(common.ErrEmpty, onDetails+": len(item.DetailsRaw) < 1")
	}

	err := bson.Unmarshal(item.DetailsRaw, exemplar)
	if err != nil {
		return errors.Wrapf(err, onDetails+": on bson.Unmarshal(%#v, %#v)", item.DetailsRaw, exemplar)
	}

	return nil
}

const onCount = "on dataMongoDB.Count()"

func (mgoOp dataMongoDB) Count(*selectors.Term, *crud.GetOptions) ([]crud.Counter, error) {
	return nil, common.ErrNotImplemented
}

const onList = "on dataMongoDB.List()"

func (mgoOp dataMongoDB) List(*selectors.Term, *crud.GetOptions) ([]data.Item, error) {

	// TODO!!!
	filter := bson.M{}

	res, err := mgoOp.collection.Find(nil, filter)

	if err != nil {
		return nil, errors.Wrapf(err, onList+": can't do collection.Find(nil, %#v)", filter)
	}
	if res == nil {
		return nil, errors.Errorf(onList+": no error and no cursor are returned with collection.Find(nil, %#v)", filter)
	}

	var briefs []data.Item

	for res.Next(nil) {
		var brief data.Item
		err := res.Decode(&brief)
		if err != nil {
			var getID mgolib.GetID
			errID := res.Decode(&getID)
			if errID != nil {
				l.Errorf(onList+": can't get _id: %s", errID)
			}

			l.Errorf(onList+": on res.Decode document with _id = %s: %s", getID.ID.Hex(), err)
			continue
		}

		briefs = append(briefs, brief)

	}

	return briefs, nil
}

const onRemove = "on dataMongoDB.Remove()"

func (mgoOp dataMongoDB) Remove(*selectors.Term, *crud.RemoveOptions) error {
	return common.ErrNotImplemented
}

const onClose = "on dataMongoDB.Close()"

func (mgoOp dataMongoDB) Close() error {
	err := mgoOp.client.Disconnect(nil)
	if err != nil {
		return errors.Wrapf(err, onClose+": can't .client.Disconnect(nil)")
	}

	return nil
}

// cleaner -----------------------------------------------------------------------------------------------------------------

var _ crud.Cleaner = &dataMongoDB{}

const onClean = "on dataMongoDB.Clean()"

func (mgoOp dataMongoDB) Clean() error {
	//filter := bson.M{"user_id": string(userID)}

	filter := bson.D{}

	num, err := mgoOp.collection.CountDocuments(nil, filter)
	if err != nil {
		return errors.Wrapf(err, onClean+": can't .CountDocuments(nil, %#v)", filter)
	}

	if num > 0 {

		l.Infof("documents to clean from collection: %d", num)

		_, err := mgoOp.collection.DeleteMany(nil, filter)
		if err != nil {
			return errors.Wrapf(err, onClean+": can't .DeleteMany(nil, %#v)", filter)
		}
	}

	return nil
}
