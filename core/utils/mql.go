package utils

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"primrose/clients"
)

var ReturnNewOption = options.FindOneAndUpdate().SetReturnDocument(options.After)
var ReturnNewReplacedOption = options.FindOneAndReplace().SetReturnDocument(options.After)
var ReturnUpsertOption = options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
var ReturnUpsertWithReplaceOption = options.FindOneAndReplace().SetUpsert(true).SetReturnDocument(options.After)

type TransactionalSession = func(ctx mongo.SessionContext) (interface{}, error)

func UseSession(a TransactionalSession) (interface{}, error) {
	txOptions := options.Transaction().SetWriteConcern(writeconcern.New(writeconcern.WMajority()))
	session, err := clients.MongoClient.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(context.TODO())
	b, e := session.WithTransaction(context.TODO(), a, txOptions)
	return b, e
}

func FromSession[T any](a TransactionalSession) (*T, error) {
	b, e := UseSession(a)
	if e != nil {
		return nil, e
	}
	r := b.(T)
	return &r, nil
}

func ReturningOne[T any](a func(t *T) error) (*T, error) {
	res, err := Returnable(a)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return res, nil
}

func Returnable[T any](a func(t *T) error) (*T, error) {
	var temp T
	err := a(&temp)
	if err != nil {
		return nil, err
	}
	return &temp, nil
}
