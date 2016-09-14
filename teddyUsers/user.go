package teddyUsers

import (
	"github.com/kataras/iris"
	"github.com/asaskevich/govalidator"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
//	"os/user"
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
	Name        string `json:"name" valid:"required" form:"name"`
	Email string `json:"email" valid:"required,email" form:"email"`
	Password string `json:"password" valid:"required" form:"password"`
	CheckPassword string `json:"password2" valid:"required" form:"password2"`
}


func (storage *UserStorage) CreateUser(ctx *iris.Context) (*NewUser, error){
	ctx.Request.Body()
	var user = NewUser{}
	err := ctx.ReadForm(&user)
	if err !=nil{
		return nil, err
	}
	err = user.check()
	if err != nil{
		return nil, err
	}
	hashpassword,err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	err = storage.Engine.Create(user.Name, user.Email, string(hashpassword))
	if err != nil{
		return nil, err
	}
	ctx.Session().Set("name", user.Name)
	return &user, nil
}

func (storage *UserStorage) CreateUserApi(ctx *iris.Context) (*NewUser, error){
	ctx.Request.Body()
	var user NewUser
	err := json.Unmarshal(ctx.Request.Body(), &user)
	if err !=nil{
		return nil, err
	}

	err = user.check()
	if err != nil{
		return nil, err
	}

	hashpassword,err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	err = storage.Engine.Create(user.Name, user.Email, string(hashpassword))
	if err != nil{
		return nil, err
	}
	ctx.Session().Set("name", user.Name)
	return &user, nil
}

func (user *NewUser) check() error {
	if check, err := govalidator.ValidateStruct(user); !check{
		return err
	} else if user.Password != user.CheckPassword{
		return errors.New("passwords mismatch")
	}
	return nil
}

func New(cookiepath string) *UserStorage{
	config := Config{
		SessionCookieName: cookiepath,
	}
	return &UserStorage{Config:&config}
}