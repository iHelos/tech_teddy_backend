package user

import (
	"github.com/kataras/iris"

	"encoding/json"
	"strings"
	"github.com/tarantool/go-tarantool"
	"github.com/iHelos/tech_teddy/helper/REST"
)

type UserStorageEngine interface {
	Create(string, string, string) (error)
	Load(string) (NewUser, error)
	CheckLogin(string, string) (error)
	CheckIsLogged(string) (error)
}

type UserStorage struct {
	Config *Config
	Engine UserStorageEngine
}

func (storage *UserStorage) CreateUser(ctx *iris.Context) (error) {
	ctx.Request.Body()
	var user = NewUser{}
	var err error
	err = json.Unmarshal(ctx.Request.Body(), &user)
	if err != nil {
		UserError := NewUserError()
		UserError.Append("request", 0)
		return UserError
	}
	user.Login = strings.ToLower(user.Login)
	err = user.check()
	if err != nil {
		return err
	}
	err = storage.Engine.Create(user.Login, user.Email, user.Password1)
	if err != nil {
		UserError := NewUserError()
		if trnerror, ok := err.(tarantool.Error); ok {
			if trnerror.Code == 32771 {
				UserError.Append("login", 1)
				return UserError
			}
		}
		UserError.Append("DB", 0)
		return UserError
	}
	ctx.Session().Set("name", user.Login)
	return nil
}

type LoginUser struct {
	Login    string `json:"name" valid:"required" form:"name"`
	Password string `json:"password" valid:"required" form:"password"`
}

func (storage *UserStorage) LoginUser(ctx *iris.Context) (error) {
	ctx.Request.Body()
	var user LoginUser
	var err error
	err = json.Unmarshal(ctx.Request.Body(), &user)
	if err != nil {
		UserError := NewUserError()
		UserError.Append("request", 0)
		return UserError
	}
	user.Login = strings.ToLower(user.Login)
	err = storage.Engine.CheckLogin(user.Login, user.Password)
	if err != nil {
		UserError := NewUserError()
		UserError.Append("request", 1)
		return UserError
	}
	ctx.Session().Set("name", user.Login)
	return nil
}

func New(cookiepath string) *UserStorage {
	config := Config{
		SessionCookieName: cookiepath,
	}
	return &UserStorage{Config:&config}
}

func (storage *UserStorage) MustBeLogged(ctx *iris.Context) {
	sid := ctx.GetCookie(storage.Config.SessionCookieName)
	err := storage.Engine.CheckIsLogged(sid)
	if err != nil {
		ctx.JSON(iris.StatusOK, REST.GetResponse(
			1,
			map[string]int{
				"loginstatus":0,
			},
		))
	} else {
		ctx.Next()
	}
}