package mimi

import (
	"context"
	"github.com/kataras/golog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"primrose/clients"
	"time"
)

var AcquiredLocks = make(map[string]bool)

func Lock(context context.Context, resource string) (bool, error) {
	_, err := clients.MongoClient.Database(Database).
		Collection("locks").
		InsertOne(context, bson.M{"_id": resource})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return false, nil
		}
		return false, err
	}
	AcquiredLocks[resource] = true
	return true, nil
}

func Release(context context.Context, resource string) (bool, error) {
	_, err := clients.MongoClient.Database(Database).
		Collection("locks").
		DeleteOne(context, bson.M{"_id": resource})
	if err != nil {
		return false, err
	}
	delete(AcquiredLocks, resource)
	return true, nil
}

func ReleaseAllLocks(ctx context.Context) {
	golog.Warn("[DB] Releasing all acquired locks.")
	var locks []string
	for k := range AcquiredLocks {
		locks = append(locks, k)
	}
	if len(locks) > 0 {
		_, err := clients.MongoClient.Database(Database).
			Collection("locks").
			DeleteMany(ctx, bson.M{"_id": bson.M{"$in": locks}})
		if err != nil {
			golog.Error("[DB] Failed to release all acquired locks, manual action is needed: ", err)
		}
		AcquiredLocks = map[string]bool{}
	}
}

func WithLock(resource string, action func()) {
	go func() {
		golog.Info("[DB] Attempting to lock ", resource, ".")
		locked, err := Lock(context.TODO(), resource)
		if err != nil {
			golog.Fatal("[DB] Failed to acquire lock for ", resource, ": ", err)
			return
		}
		if !locked {
			golog.Warn("[DB] Lock ", resource, " is already acquired by another resource, giving up.")
			return
		}
		golog.Info("[DB] Acquired ", resource, " lock, going ahead with function.")
		defer func(context context.Context, resource string) {
			time.Sleep(3 * time.Second)
			_, err := Release(context, resource)
			if err != nil {
				golog.Error("[DB] Failed to release lock for ", resource, ": ", err)
				return
			}
			golog.Info("[DB] Released lock for ", resource, ".")
		}(context.TODO(), resource)
		action()
	}()
}

func TryLock(resource string, action func()) {
	go func() {
		golog.Info("[DB] Attempting to lock ", resource, ".")
		locked, err := Lock(context.TODO(), resource)
		if err != nil {
			golog.Fatal("[DB] Failed to acquire lock for ", resource, ": ", err)
			return
		}
		if !locked {
			golog.Warn("[DB] Lock ", resource, " is already acquired by another resource, giving up.")
			return
		}
		golog.Info("[DB] Acquired ", resource, " lock, going ahead with function.")
		action()
	}()
}
