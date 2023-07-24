package comments

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"primrose/models"
	"primrose/models/users"
	"primrose/utils"
	"time"
)

const (
	MaximumContentLength = 5_120
)

type Comment struct {
	models.Model `bson:",inline"`
	Content      string              `json:"content" bson:"content" validate:"required"`
	Post         string              `json:"post" bson:"post"`
	ReplyTo      *primitive.ObjectID `json:"reply_to" bson:"reply_to"`
	CreatedAt    time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at" bson:"updated_at"`
}

type UnjoinedComment struct {
	Comment `bson:",inline"`
	Author  primitive.ObjectID `json:"author" bson:"author"`
}

type JoinedComment struct {
	Comment `bson:",inline"`
	Author  users.User `json:"author" bson:"author"`
}

func (joinedComment *JoinedComment) Unjoin() *UnjoinedComment {
	return &UnjoinedComment{Comment: joinedComment.Comment, Author: joinedComment.Author.ObjectId()}
}

func (comment *Comment) Delete() error {
	if _, err := GetCollection().DeleteMany(context.TODO(), bson.M{"reply_to": comment.ObjectId()}); err != nil {
		return err
	}
	if _, err := GetCollection().DeleteOne(context.TODO(), bson.M{"_id": comment.ObjectId()}); err != nil {
		return err
	}
	return nil
}

func (comment *UnjoinedComment) Parent() (*UnjoinedComment, error) {
	if comment.ReplyTo == nil {
		return nil, nil
	}
	return WithId(*comment.ReplyTo)
}

func (comment *UnjoinedComment) User() (*users.User, error) {
	return users.WithId(comment.Author)
}

var (
	CannotReplyToUnknownCommentErr = errors.New("cannot reply to unknown comment")
)

func (comment *UnjoinedComment) Save() error {
	if comment.Id == "" {
		if comment.ReplyTo != nil {
			parent, err := comment.Parent()
			if err != nil {
				return err
			}
			if parent == nil {
				return CannotReplyToUnknownCommentErr
			}
		}
		author, err := comment.User()
		if err != nil {
			return err
		}
		if author == nil {
			return errors.New("cannot find the author of the comment")
		}
	}
	comment.UpdatedAt = time.Now()
	if comment.Id == "" {
		result, err := GetCollection().InsertOne(context.TODO(), comment)
		if err != nil {
			return err
		}
		comment.Id = result.InsertedID.(primitive.ObjectID).Hex()
		return nil
	}
	return GetCollection().FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": comment.ObjectId()},
		bson.M{"$set": bson.M{
			"updated_at": comment.UpdatedAt,
			"content":    comment.Content,
		}},
		utils.ReturnUpsertOption,
	).Decode(&comment)
}
