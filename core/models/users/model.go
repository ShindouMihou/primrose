package users

import (
	"primrose/models"
	"primrose/utils"
)

type User struct {
	models.Model
	Name     string `json:"name" bson:"name"`
	Email    string `json:"-" bson:"email"`
	Password string `json:"-" bson:"password"`
	Token    string `json:"-" bson:"token"`
	Flags    []Flag `json:"flags" bson:"flags"`
}

type Flag = string

const (
	AdminFlag Flag = "admin"
)

func (user *User) IsAdmin() bool {
	return utils.AnyMatchString(user.Flags, AdminFlag)
}
