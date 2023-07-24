package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Model struct {
	Id string `json:"id" bson:"_id,omitempty"`
}

func (model *Model) ObjectId() primitive.ObjectID {
	oid, _ := primitive.ObjectIDFromHex(model.Id)
	return oid
}
