package posts

import (
	"go.mongodb.org/mongo-driver/mongo"
	"primrose/clients"
)

func GetCollection() *mongo.Collection {
	return clients.Db.Collection("posts")
}
