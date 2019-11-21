package crud_mongodb

import (
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"reflect"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/libraries/mgolib"
	"github.com/pavlo67/workshop/components/crud"
	"github.com/pavlo67/workshop/components/selector"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO: use common context.Session!!!

// storage -----------------------------------------------------------------------------------------------------------------

var _ crud.Operator = &crudMongoDB{}

type crudMongoDB struct {
	client     *mongo.Client
	collection *mongo.Collection
	exemplar   interface{}
}

const onNewCRUD = "on crud_mongodb.NewCRUD()"

func NewCRUD(dbAccess *config.Access, timeout time.Duration, collectionName string, exemplar interface{}) (crud.Operator, crud.Cleaner, *mongo.Client, error) {
	if exemplar == nil {
		return nil, nil, nil, errors.New("no exemplar")
	}

	client, err := mgolib.Connect(dbAccess, timeout)
	if err != nil {
		return nil, nil, nil, err
	}
	database := client.Database(dbAccess.Path)
	mgoOp := &crudMongoDB{
		client:     client,
		collection: database.Collection(collectionName),
		exemplar:   exemplar,
	}

	return mgoOp, mgoOp, client, nil
}

// operator ----------------------------------------------------------------------------------------------------------------

const onExemplar = "on crudMongoDB.Exemplar()"

func (mgoOp crudMongoDB) Exemplar() interface{} {
	if reflect.TypeOf(mgoOp.exemplar).Kind() == reflect.Ptr {
		// Pointer:
		return reflect.New(reflect.ValueOf(mgoOp.exemplar).Elem().Type()).Interface()
	}

	// Not pointer:
	return reflect.New(reflect.TypeOf(mgoOp.exemplar)).Elem().Interface()
}

const onSave = "on crudMongoDB.Save()"

func (mgoOp crudMongoDB) Save(item crud.Item, options *crud.SaveOptions) (*common.ID, error) {
	// TODO: use Upsert

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

	res, err := mgoOp.collection.InsertOne(nil, &item)

	if err != nil {
		return nil, errors.Wrapf(err, onSave+": can't .InsertOne(nil, %#v)", item)
	}

	var id common.ID
	if objectID, ok := res.InsertedID.(primitive.ObjectID); ok {
		id = common.ID(objectID.Hex())
	} else {
		l.Debugf("res.InsertedID = %#v", res.InsertedID)
	}

	return &id, nil
}

const onRead = "on crudMongoDB.Read()"

func (mgoOp crudMongoDB) Read(id common.ID, options *crud.GetOptions) (*crud.Item, error) {
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

	item := crud.Item{Details: mgoOp.Exemplar()}

	err = res.Decode(&item)
	if err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			return nil, nil
		}

		return nil, errors.Wrapf(err, onRead+": on res.Decode(%#v)", res)
	}

	return &item, nil
}

const onExists = "on crudMongoDB.Exists()"

func (mgoOp crudMongoDB) Exists(*selector.Term, *crud.GetOptions) ([]crud.Part, error) {
	return nil, common.ErrNotImplemented
}

const onList = "on crudMongoDB.List()"

func (mgoOp crudMongoDB) List(*selector.Term, *crud.GetOptions) ([]crud.Brief, error) {

	// TODO!!!
	filter := bson.M{}

	res, err := mgoOp.collection.Find(nil, filter)

	if err != nil {
		return nil, errors.Wrapf(err, onList+": can't do collection.Find(nil, %#v)", filter)
	}
	if res == nil {
		return nil, errors.Errorf(onList+": no error and no cursor are returned with collection.Find(nil, %#v)", filter)
	}

	var briefs []crud.Brief

	for res.Next(nil) {
		var brief crud.Brief
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

const onRemove = "on crudMongoDB.Remove()"

func (mgoOp crudMongoDB) Remove(*selector.Term, *crud.RemoveOptions) error {
	return common.ErrNotImplemented
}

const onClose = "on crudMongoDB.Close()"

func (mgoOp crudMongoDB) Close() error {
	err := mgoOp.client.Disconnect(nil)
	if err != nil {
		return errors.Wrapf(err, onClose+": can't .client.Disconnect(nil)")
	}

	return nil
}

// cleaner -----------------------------------------------------------------------------------------------------------------

var _ crud.Cleaner = &crudMongoDB{}

const onClean = "on crudMongoDB.Clean()"

func (mgoOp crudMongoDB) Clean() error {
	//filter := bson.M{"user_id": string(userID)}

	_, err := mgoOp.collection.DeleteMany(nil, nil)
	if err != nil {
		return errors.Wrapf(err, onClean+": can't .DeleteMany(nil, %#v, nil)", nil)
	}

	return nil
}
