package toolkit

import (
	"github.com/kataras/golog"
	"os"
	"primrose/clients"
	"primrose/toolkit/collections"
	"primrose/toolkit/gen"
	"primrose/toolkit/indexes"
	"primrose/toolkit/migrations"
	"primrose/toolkit/mimi"
	"strings"
)

func HandleArguments() bool {
	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case "--migrate":
			migrations.Discover()
			params := args[1:]
			if len(params) > 0 {
				mimi.MigrateOne(params[0], clients.MongoClient)
				return false
			}

			mimi.Migrate(clients.MongoClient)
			return false
		case "--rollback":
			migrations.Discover()
			params := args[1:]
			if len(params) > 0 {
				mimi.RollbackOne(params[0], clients.MongoClient)
				return false
			}

			mimi.Rollback(clients.MongoClient)
			return false
		case "--create-migration":
			params := args[1:]
			if len(params) < 2 {
				golog.Fatal("Required options: [collection] [...name]")
			}
			name := strings.Join(params[1:], " ")
			gen.WriteMigration(name, params[0])
			return false
		case "--create-index":
			params := args[1:]
			if len(params) < 2 {
				golog.Fatal("Required options: [collection] [...name]")
			}
			name := strings.Join(params[1:], " ")
			gen.WriteIndex(name, params[0])
			return false
		case "--index":
			indexes.Discover()
			params := args[1:]
			if len(params) > 0 {
				mimi.IndexOne(params[0], clients.MongoClient)
				return false
			}

			mimi.IndexAll(clients.MongoClient)
			return false
		case "--create-collection":
			params := args[1:]
			if len(params) < 1 {
				golog.Fatal("Required options: [...name]")
			}
			name := strings.Join(params, " ")
			gen.WriteCollection(name)
			return false
		case "--migrate-collections":
			collections.Discover()
			params := args[1:]
			if len(params) > 0 {
				mimi.CreateCollection(params[0], clients.MongoClient)
				return false
			}

			mimi.CreateAllCollections(clients.MongoClient)
			return false
		}
	}
	return true
}
