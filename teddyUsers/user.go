package teddyUsers

import (
	"github.com/kataras/iris"
	"github.com/asaskevich/govalidator"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type UserStorageEngine interface {
	Create(string, string, string) (error)
	Load(string) (NewUser, error)
}

type UserStorage struct {
	Config *Config
	Engine UserStorageEngine
}

type NewUser struct {
	Name        string `json:"name" valid:"required"`
	Email string `json:"email" valid:"required,email"`
	Password string `json:"password" valid:"required"`
	CheckPassword string `json:"password2" valid:"required"`
}

func (storage *UserStorage) CreateUser(ctx *iris.Context) (*NewUser, error){
	ctx.Request.Body()
	var user NewUser
	err := json.Unmarshal(ctx.Request.Body(), &user)
	if err !=nil{
		return nil, err
	}

	if check, err := govalidator.ValidateStruct(user); !check{
		return nil, err
	}

	if user.Password != user.CheckPassword{
		return nil, errors.New("passwords mismatch")
	}

	hashpassword,err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	storage.Engine.Create(user.Name, user.Email, string(hashpassword))
	return &user, nil
}

func New(cookiepath string) *UserStorage{
	config := Config{
		SessionCookieName: cookiepath,
	}
	return &UserStorage{Config:&config}
}