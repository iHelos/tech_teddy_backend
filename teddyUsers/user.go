package teddyUsers

import (
	"github.com/kataras/iris"
	"github.com/asaskevich/govalidator"
	"encoding/json"
	"errors"
	//	"os/user"
)

type UserStorageEngine interface {
	Create(string, string, string) (error)
	Load(string) (NewUser, error)
	CheckLogin(string, string) (error)
}

type UserStorage struct {
	Config *Config
	Engine UserStorageEngine
}

type NewUser struct {
	Name          string `json:"name" valid:"required" form:"name"`
	Email         string `json:"email" valid:"required,email" form:"email"`
	Password      string `json:"password" valid:"required" form:"password"`
	CheckPassword string `json:"password2" valid:"required" form:"password2"`
}

func (storage *UserStorage) CreateUser(ctx *iris.Context) (error) {
	ctx.Request.Body()
	var user = NewUser{}
	var err error
	if string(ctx.Request.Header.ContentType()) == "application/json" {
		err = json.Unmarshal(ctx.Request.Body(), &user)
	} else {
		err = ctx.ReadForm(&user)
	}
	if err != nil {
		return err
	}
	err = user.check()
	if err != nil {
		return err
	}
	err = storage.Engine.Create(user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}
	ctx.Session().Set("name", user.Name)
	return nil
}

type LoginUser struct {
	Name     string `json:"name" valid:"required" form:"name"`
	Password string `json:"password" valid:"required" form:"password"`
}

func (storage *UserStorage) LoginUser(ctx *iris.Context) (error) {
	ctx.Request.Body()
	var user LoginUser
	var err error
	if string(ctx.Request.Header.ContentType()) == "application/json" {
		err = json.Unmarshal(ctx.Request.Body(), &user)
	} else {
		err = ctx.ReadForm(&user)
	}
	if err != nil {
		return err
	}

	err = user.check()
	if err != nil {
		return err
	}

	err = storage.Engine.CheckLogin(user.Name, user.Password)
	if err != nil {
		return err
	}
	ctx.Session().Set("name", user.Name)
	return nil
}

func (user *NewUser) check() error {
	if check, err := govalidator.ValidateStruct(user); !check {
		return err
	} else if user.Password != user.CheckPassword {
		return errors.New("passwords mismatch")
	}
	return nil
}

func (user *LoginUser) check() error {
	if check, err := govalidator.ValidateStruct(user); !check {
		return err
	}
	return nil
}

func New(cookiepath string) *UserStorage {
	config := Config{
		SessionCookieName: cookiepath,
	}
	return &UserStorage{Config:&config}
}