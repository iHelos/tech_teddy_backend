package user

import (
	"github.com/kataras/iris"

	"encoding/json"
	"strings"
	"github.com/tarantool/go-tarantool"
	"github.com/iHelos/tech_teddy/helper/REST"
	"github.com/streamrail/concurrent-map"
	"github.com/NaySoftware/go-fcm"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/dgrijalva/jwt-go"
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

var fcmList cmap.ConcurrentMap

const (
	hmacUserSecret = "95CCEB5921E59B285AC773E4963E1"
	serverKey = "AIzaSyACAJnKTfYG9_gmBna5TU-VA57ssTl3TVk"
)

func (storage *UserStorage) CreateUser(ctx *iris.Context) (string,string,error) {
	log.Print(ctx.Request.Header.String())
	var user = NewUser{}
	var err error
	err = json.Unmarshal(ctx.Request.Body(), &user)
	if err != nil {
		UserError := NewUserError()
		UserError.Append("request", 0)
		return "","",UserError
	}
	user.Login = strings.ToLower(user.Login)
	err = user.check()
	if err != nil {
		return "","",err
	}
	err = storage.Engine.Create(user.Login, user.Email, user.Password1)
	if err != nil {
		UserError := NewUserError()
		if trnerror, ok := err.(tarantool.Error); ok {
			if trnerror.Code == 32771 {
				UserError.Append("login", 1)
				return "","", UserError
			}
		}
		UserError.Append("DB", 0)
		return "","", UserError
	}
	//ctx.Session().Set("name", user.Login)
	userToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": user.Login,
		"type":"user",
	})
	bearToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": user.Login,
		"type":"bear",
	})
	userTokenSigned, err := userToken.SignedString([]byte(hmacUserSecret))
	bearTokenSigned, err := bearToken.SignedString([]byte(hmacUserSecret))

	return userTokenSigned,bearTokenSigned,nil
}

type LoginUser struct {
	Login    string `json:"name" valid:"required" form:"name"`
	Password string `json:"password" valid:"required" form:"password"`
	FCMToken string `json:"fcm"`
}

func (storage *UserStorage) LoginUser(ctx *iris.Context) (string,string,error) {
	var user LoginUser
	var err error
	err = json.Unmarshal(ctx.Request.Body(), &user)
	if err != nil {
		UserError := NewUserError()
		UserError.Append("request", 0)
		return "","", UserError
	}
	user.Login = strings.ToLower(user.Login)
	err = storage.Engine.CheckLogin(user.Login, user.Password)
	if err != nil {
		UserError := NewUserError()
		UserError.Append("request", 1)
		return "","",UserError
	}
	userToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": user.Login,
		"type":"user",
	})
	bearToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": user.Login,
		"type":"bear",
	})
	userTokenSigned, err := userToken.SignedString([]byte(hmacUserSecret))
	bearTokenSigned, err := bearToken.SignedString([]byte(hmacUserSecret))

	return userTokenSigned,bearTokenSigned,nil
}

func New(cookiepath string) *UserStorage {
	config := Config{
		SessionCookieName: cookiepath,
	}
	return &UserStorage{Config:&config}
}

func (storage *UserStorage) MustBeLogged(ctx *iris.Context) {

	signedtoken := ctx.Request.Header.Peek("Authorization")
	token, err := jwt.Parse(string(signedtoken), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(hmacUserSecret), nil
	})
	if (err!=nil){
		ctx.JSON(iris.StatusOK, REST.GetResponse(
			1,
			map[string]int{
				"loginstatus":0,
			},
		))
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["login"])
	} else {
		fmt.Println(err)
	}
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

func SendAll() {
	data := map[string]string{
		"msg": "Что за странный медведь",
		"sum": "ГОВОРИ!",
	}
	//fcmList.Set("dIJBF6MaTCo:APA91bG1FB4K37boqIj-E3rv3KUTjopWh6sXa5IcBhWUfOw9mSFcmiyHgY4eAKSyOPWj1cGpNZJrDkkTGsO2Wbzb0J59rvddn1Kn2PnRn4o2C9miS9QQbPAabrNwM8wFLXamC26T37ZQ", "")
	//fcmList.Set("c8FnS_Yc478:APA91bHoh6yfYQY9id09ZFJkD3sKuBI7VqBJPQActJ1Ra9QoXDowMKyZ0fdKe2mRIrD11YSVCH-1Kfv0TVtFdmAsl6bjEFJQOhkN-3hfgoTYSXW3grCGLEsvT2vv_-Y0weoLBo8VOPjw", "")

	log.Print("sending to...")
	c := fcm.NewFcmClient(serverKey)
	devices := fcmList.Keys()
	log.Print("1) ", devices[0])
	log.Print("2) ", devices[1])
	if len(devices) > 1 {
		c.NewFcmMsgTo("", data)
		c.AppendDevices(devices)
		notific := fcm.NotificationPayload{
			Title: "Мишка",
			Body: "ПРИВЕТ!",
		}
		c.SetNotificationPayload(&notific)

		status, err := c.Send()
		if err == nil {
			status.PrintResults()
		} else {
			fmt.Println(err)
		}
	}
}