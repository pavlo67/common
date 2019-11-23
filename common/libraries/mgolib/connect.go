package mgolib

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"log"

	"github.com/pavlo67/workshop/common/config"
)

type GetID struct {
	ID primitive.ObjectID `json:"id,omitempty"   bson:"_id,omitempty"`
}

func Address(access *config.Access) string {
	if access == nil {
		return ""
	}

	return fmt.Sprintf("mongodb://%s:%d/", access.Host, access.Port)
}

func Connect(access *config.Access, timeout time.Duration) (*mongo.Client, error) {
	if access == nil {
		return nil, errors.New("no config data for MongoDB")
	}

	clientOptions := options.Client().ApplyURI(Address(access))

	if access.User != "" {
		clientOptions.SetAuth(options.Credential{
			// AuthMechanism:           "",
			// AuthMechanismProperties: nil,
			// PasswordSet:             false,

			AuthSource: "admin", // TODO!!!
			Username:   access.User,
			Password:   access.Pass,
		})
	}

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, errors.Wrapf(err, "can't mongo.NewClient(options.Client().ApplyURI(Address(%#v)))", access)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "can't client.Connect(%#v)", access)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "can't client.Ping(%#v)", access)
	}

	log.Printf("connected to %#v", access)

	return client, nil
}
