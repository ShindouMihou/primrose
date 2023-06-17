package mimi

import (
	"context"
	"errors"
	"github.com/kataras/golog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"primrose/utils"
	"time"
)

var Indexes []Index
var Migrations []Migration
var Collections []Collection

var Database = "primrose"

func AddMigration(migration Migration) *Migration {
	Migrations = append(Migrations, migration)
	return &migration
}

func AddCollection(collection Collection) *Collection {
	Collections = append(Collections, collection)
	return &collection
}

func AddIndex(index Index) *Index {
	Indexes = append(Indexes, index)
	return &index
}

func Migrate(conn *mongo.Client) {
	golog.Info("[DB] Ignoring all completed migrations")
	completedMigrations := make(map[string]iMigration)
	cMigrations := GetCompletedMigrations(conn)
	for _, m := range cMigrations {
		completedMigrations[m.Id] = m
	}
	golog.Info("[DB] Preparing to migrate ", len(Migrations), " migrations into the database.")
	_, err := utils.UseSession(func(ctx mongo.SessionContext) (interface{}, error) {
		for _, value := range Migrations {
			migration := value
			if _, completed := completedMigrations[migration.Key]; completed {
				golog.Warn("[DB] Migration ", migration.Key, " has already been completed, therefore, it is skipped.")
				continue
			}
			if err := ExecMigration(ctx, migration, conn); err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
	if err != nil {
		golog.Fatal("[DB] Failed to complete migration due to error, no changes have occurred: ", err)
		return
	}
	golog.Info("[DB] Successfully migrated ", len(Migrations), " migrations")
}

func GetMigrationsCollection(conn *mongo.Client) *mongo.Collection {
	return conn.Database(Database).Collection("migrations")
}

func createMigrationResult(migration *Migration) iMigration {
	return iMigration{Id: migration.Key, CreatedAt: time.Now()}
}

func ExecMigration(ctx mongo.SessionContext, migration Migration, conn *mongo.Client) error {
	if _, err := GetMigrationsCollection(conn).InsertOne(ctx, createMigrationResult(&migration)); err != nil {
		return errors.Join(errors.New("cannot insert migration result"), err)
	}
	return migration.MigrateUp(ctx, conn.Database(Database).Collection(migration.Collection))
}

func FindMigration(key string) *Migration {
	migration := (*Migration)(nil)
	for _, m := range Migrations {
		if m.Key == key {
			migration = &m
			break
		}
	}
	return migration
}

func MigrateOne(key string, conn *mongo.Client) {
	migration := FindMigration(key)
	if migration == nil {
		golog.Fatal("[DB] Cannot find the migration ", key)
	}
	if migration.MigrateUp != nil {
		golog.Info("[DB] Preparing to migrate ", key, " into the database.")
		_, err := utils.UseSession(func(ctx mongo.SessionContext) (interface{}, error) {
			return nil, migration.MigrateUp(ctx, conn.Database(Database).Collection(migration.Collection))
		})
		if err != nil {
			golog.Fatal("[DB] Failed to migrate ", key, " due to error: ", err)
			return
		}
	}
}

func ExecRollback(ctx mongo.SessionContext, migration iMigration, conn *mongo.Client) error {
	migrationsCollection := GetMigrationsCollection(conn)
	m := FindMigration(migration.Id)
	if m == nil {
		return errors.New("cannot find rollback plan for " + m.Key)
	}
	_, err := migrationsCollection.DeleteOne(ctx, bson.M{"_id": migration.Id})
	if err != nil {
		return errors.Join(errors.New("cannot release migration "+m.Key), err)
	}
	err = m.MigrateDown(ctx, conn.Database(Database).Collection(m.Collection))
	if err != nil {
		return errors.Join(errors.New("cannot rollback migration "+m.Key), err)
	}
	return nil
}

func RollbackOne(key string, conn *mongo.Client) {
	var migration iMigration

	migrationsCollection := GetMigrationsCollection(conn)
	err := migrationsCollection.FindOne(context.TODO(), bson.M{"_id": key}).Decode(&migration)
	if err != nil {
		golog.Fatal("[DB] Cannot rollback ", key, " due to not being able to get record of it: ", err)
	}
	_, err = utils.UseSession(func(ctx mongo.SessionContext) (interface{}, error) {
		return nil, ExecRollback(ctx, migration, conn)
	})
	if err != nil {
		golog.Fatal("[DB] Cannot rollback ", key, ": ", err)
	}
	golog.Info("[DB] Successfully rolled back ", key)
}

func GetCompletedMigrations(conn *mongo.Client) []iMigration {
	var migrations []iMigration
	migrationsCollection := GetMigrationsCollection(conn)
	cursor, err := migrationsCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []iMigration{}
		}
		golog.Fatal("[DB] Cannot retrieve list of migrations: ", err)
	}

	if err = cursor.All(context.TODO(), &migrations); err != nil {
		golog.Fatal("[DB] Cannot retrieve list of migrations: ", err)
	}
	return migrations
}

func Rollback(conn *mongo.Client) {
	migrations := GetCompletedMigrations(conn)
	_, err := utils.UseSession(func(ctx mongo.SessionContext) (interface{}, error) {
		for _, migration := range migrations {
			return nil, ExecRollback(ctx, migration, conn)
		}
		return nil, nil
	})
	if err != nil {
		golog.Fatal("[DB] Cannot rollback migrations: ", err)
	}
	golog.Info("[DB] Successfully rolled back ", len(migrations), " migrations.")
}

func FindIndex(key string) *Index {
	index := (*Index)(nil)
	for _, i := range Indexes {
		if i.Key == key {
			index = &i
			break
		}
	}
	return index
}

func ExecIndex(index Index, conn *mongo.Client) {
	if err := index.Action(conn.Database(Database).Collection(index.Collection)); err != nil {
		golog.Error("[DB] Failed to complete index migration ", index.Key, ": ", err)
		return
	}
	golog.Info("[DB] Successfully completed index migration ", index.Key)
}

func IndexAll(conn *mongo.Client) {
	for _, index := range Indexes {
		ExecIndex(index, conn)
	}
}

func IndexOne(key string, conn *mongo.Client) {
	index := FindIndex(key)
	if index == nil {
		golog.Fatal("[DB] Cannot find any index migration with the key of ", key)
	}
	ExecIndex(*index, conn)
}

func FindCollection(key string) *Collection {
	collection := (*Collection)(nil)
	for _, i := range Collections {
		if i.Key == key {
			collection = &i
			break
		}
	}
	return collection
}

func ExecCollection(collection Collection, conn *mongo.Client) {
	if err := collection.Action(conn.Database(Database)); err != nil {
		var commandError mongo.CommandError
		if errors.As(err, &commandError) {
			golog.Warn("[DB] Ignorable failure to complete collection migration ", collection.Key, " most likely due to collection existing.")
			return
		}
		golog.Error("[DB] Failed to complete collection migration ", collection.Key, ": ", err)
		return
	}
	golog.Info("[DB] Successfully completed collection migration ", collection.Key)
}

func CreateAllCollections(conn *mongo.Client) {
	for _, collection := range Collections {
		ExecCollection(collection, conn)
	}
}

func CreateCollection(key string, conn *mongo.Client) {
	collection := FindCollection(key)
	if collection == nil {
		golog.Fatal("[DB] Cannot find any collection migration with the key of ", key)
	}
	ExecCollection(*collection, conn)
}
