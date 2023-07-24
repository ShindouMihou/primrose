package posts

import (
	"context"
	"github.com/dchest/uniuri"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"primrose/models/comments"
	"primrose/utils"
	"regexp"
	"strings"
	"time"
)

type Post struct {
	Id        string    `json:"id" bson:"_id"`
	Image     *string   `json:"image" bson:"image" validate:"http_url"`
	Title     string    `json:"title" bson:"title" validate:"required"`
	Content   string    `json:"content" bson:"content"`
	Published bool      `json:"published" bson:"published" validate:"boolean"`
	Slug      string    `json:"slug" bson:"slug"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

var SymbolRegex = regexp.MustCompile("[^A-Za-z0-9]+")

func (post *Post) Transfer(res *Post) {
	res.Image = post.Image
	if post.Title != "" {
		res.Title = post.Title
	}
	if post.Content != "" {
		res.Content = post.Content
	}
	res.Published = post.Published
	res.Slug = post.Slug
}

func (post *Post) Save() error {
	newId := false
	if post.Id == "" {
		// Primrose is created for personal blogs, which tend to have way less than 9,999 posts, therefore, this is the
		// best result.
		post.Id = uniuri.NewLen(4)
		newId = true
	}
	if post.Slug == "" {
		post.Slug = strings.ReplaceAll(post.Title, " ", "-")
		post.Slug = SymbolRegex.ReplaceAllString(post.Title, "-")
		post.Slug = strings.ToLower(post.Slug)
	}
	post.UpdatedAt = time.Now()
	if err := GetCollection().FindOneAndReplace(
		context.TODO(),
		bson.M{"_id": post.Id},
		post,
		utils.ReturnUpsertWithReplaceOption,
	).Err(); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			if newId == true {
				post.Id = uniuri.NewLen(4)
			} else {
				post.Slug = post.Slug + "-" + uniuri.NewLen(4)
			}
			return post.Save()
		}
		return err
	}
	return nil
}

func (post *Post) Delete() error {
	if post.Id == "" {
		return PostDoesNotExistErr
	}
	if _, err := comments.GetCollection().DeleteMany(context.TODO(), bson.M{"post": post.Id}); err != nil {
		return err
	}
	if _, err := GetCollection().DeleteOne(context.TODO(), bson.M{"_id": post.Id}); err != nil {
		return err
	}
	return nil
}
