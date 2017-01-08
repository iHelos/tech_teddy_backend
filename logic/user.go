package logic

import (
	"github.com/iHelos/tech_teddy/model"
	"strings"
	"github.com/kataras/iris"
	"github.com/dgrijalva/jwt-go"
	"log"
	"github.com/tarantool/go-tarantool"
	. "github.com/iHelos/tech_teddy/helper"
	"golang.org/x/crypto/bcrypt"
	"fmt"
	"github.com/kataras/go-errors"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"time"
	"math/rand"
	"crypto/md5"
	"encoding/hex"
)

const (
	hmacUserSecret = "95CCEB5921E59B285AC773E4963E1"
)

func ParseToken(ctx *iris.Context) (int, error) {
	signedtoken := ctx.Request.Header.Get("Authorization")
	log.Print(signedtoken)
	token, err := jwt.Parse(string(signedtoken), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(hmacUserSecret), nil
	})
	if (err != nil) {
		return 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if login, ok := claims["id"].(float64); ok {
			return int(login), nil
		} else {
			return 0, errors.New("bad jwt")
		}
	} else {
		return 0, err
	}
}

func Register(ctx *iris.Context) (string, string, error) {
	var user = model.NewProfile{}
	var err error
	err = ctx.ReadJSON(&user)
	if err != nil {
		log.Print(err)
		UserError := NewError()
		UserError.Append("request", 0)
		return "", "", UserError
	}
	user.Email = strings.ToLower(user.Email)
	err = user.Validate()
	if err != nil {
		return "", "", err
	}
	hashpassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password1), bcrypt.DefaultCost)
	var profile = model.Profile{Email:user.Email, Name:user.Login, Password:string(hashpassword)}
	created, err := model.CreateProfile(profile)
	if err != nil {
		UserError := NewError()
		if trnerror, ok := err.(tarantool.Error); ok {
			if trnerror.Code == 32771 || trnerror.Code == 3 {
				UserError.Append("email", 1)
				return "", "", UserError
			}
		}
		UserError.Append("DB", 0)
		return "", "", UserError
	}
	//ctx.Session().Set("name", user.Login)
	userToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  created.ID,
		"type":"user",
	})
	bearToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  created.ID,
		"type":"bear",
	})
	userTokenSigned, err := userToken.SignedString([]byte(hmacUserSecret))
	bearTokenSigned, err := bearToken.SignedString([]byte(hmacUserSecret))

	return userTokenSigned, bearTokenSigned, nil
}
func Login(ctx *iris.Context) (string, string, error) {
	var user model.Profile
	err := ctx.ReadJSON(&user)
	if err != nil {
		UserError := NewError()
		UserError.Append("request", 0)
		return "", "", UserError
	}
	user.Email = strings.ToLower(user.Email)
	id, err := model.CheckPassword(user)
	if err != nil {
		log.Print(err)
		UserError := NewError()
		UserError.Append("request", 1)
		return "", "", UserError
	}
	userToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"type":"user",
	})
	bearToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"type":"bear",
	})
	userTokenSigned, err := userToken.SignedString([]byte(hmacUserSecret))
	bearTokenSigned, err := bearToken.SignedString([]byte(hmacUserSecret))
	return userTokenSigned, bearTokenSigned, nil
}

type StoryPointer struct {
	StoryID int `json:"storyID"`
}

func LikeStory(ctx *iris.Context) (error) {
	id, err := ParseToken(ctx)
	if err != nil {
		return err
	}
	story := StoryPointer{}
	err = ctx.ReadJSON(&story)
	if err != nil {
		return err
	}
	err = model.AddStory(id, story.StoryID)
	return err
}

type vkCode struct {
	Access_token string 	`json:"access_token"`
	User_id      int        `json:"user_id"`
	Email        string        `json:"email"`
	Error        string        `json:"error"`
}


var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func VKGetCode(ctx *iris.Context) (string, string, error) {
	code := ctx.URLParam("code")
	resp, err := http.Get("https://oauth.vk.com/access_token?client_id=5806269&client_secret=QQY9VWcmlhiUrNkXXznv&redirect_uri=https://magicbackpack.ru/api/social/vk/getcode&code=" + code)
	if err != nil {
		return "", "", err

	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var answer vkCode
	json.Unmarshal(body, &answer)
	if answer.Error != "" {
		return "", "", errors.New(answer.Error)
	}
	log.Print(answer)
	profile, err := model.GetProfileEmail(answer.Email)
	if profile.Email == ""{
		profile.Email = answer.Email
		profile.Name = string(answer.User_id)
		profile.Password = randStringRunes(10)
		profile, err = model.CreateProfile(profile)
		if err != nil{
			return "", "", err
		}
	}
	userToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  profile.ID,
		"type":"user",
	})
	bearToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  profile.ID,
		"type":"bear",
	})
	userTokenSigned, err := userToken.SignedString([]byte(hmacUserSecret))
	bearTokenSigned, err := bearToken.SignedString([]byte(hmacUserSecret))
	return userTokenSigned, bearTokenSigned, nil
}

type OKResponse struct{
	Email 	string	`json:"email"`
	Name	string	`json:"name"`
	Uid 	string	`json:"uid"`
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func OKGetCode(ctx *iris.Context) (string, string, error)  {
	code := ctx.URLParam("code")
	if (code == "") {
		err_str := ctx.URLParam("error")
		return "", "", errors.New(err_str)
	}
	secret_key := "95224E00D3EE4887818C6A48"
	OKurl := fmt.Sprintf("https://api.ok.ru/oauth/token.do?code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=%s",
		code, "1249370880", secret_key, "https://magicbackpack.ru/api/social/ok/getcode", "authorization_code")

	resp, err := http.Post(OKurl,"application/json", nil)
	if err!=nil{
		return "", "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var answer vkCode
	json.Unmarshal(body, &answer)
	if answer.Error != "" {
		return "", "", errors.New(answer.Error)
	}
	md5_code := GetMD5Hash(answer.Access_token + secret_key)
	//new_code_params := string(md5_code[:]
	//params := append([]byte("application_key=CBAEOIHLEBABABABAformat=jsonmethod=users.getCurrentUser"), md5_code[:]...)
	sig := GetMD5Hash("application_key=CBAEOIHLEBABABABAformat=jsonmethod=users.getCurrentUser" + md5_code)
	new_url := fmt.Sprintf("https://api.ok.ru/fb.do?application_key=CBAEOIHLEBABABABA&format=json&method=users.getCurrentUser&sig=%s&access_token=%s", string(sig[:]), answer.Access_token)
	new_resp, err := http.Get(new_url)
	if err!=nil{
		return "", "", err
	}
	defer new_resp.Body.Close()
	body, _ = ioutil.ReadAll(new_resp.Body)
	var user_info OKResponse
	json.Unmarshal(body, &user_info)

	profile, err := model.GetProfileEmail(user_info.Uid)
	if profile.Email == ""{
		profile.Email = user_info.Uid
		profile.Name = string(user_info.Name)
		profile.Password = randStringRunes(10)
		profile, err = model.CreateProfile(profile)
		if err != nil{
			return "", "", err
		}
	}
	userToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  profile.ID,
		"type":"user",
	})
	bearToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  profile.ID,
		"type":"bear",
	})
	userTokenSigned, err := userToken.SignedString([]byte(hmacUserSecret))
	bearTokenSigned, err := bearToken.SignedString([]byte(hmacUserSecret))
	return userTokenSigned, bearTokenSigned, nil
}

type FBTokenResponse struct{
	Access_token 	string	`json:"access_token"`
	Token_type	string	`json:"token_type"`
	Expires_in 	int	`json:"expires_in"`
}

type FBUserInfo struct{
	Id 	string	`json:"id"`
	Name	string	`json:"name"`
	Email 	string	`json:"email"`
}

func FBGetCode(ctx *iris.Context) (string, string, error)  {
	code := ctx.URLParam("code")
	if (code == "") {
		err_str := ctx.URLParam("error")
		return "", "", errors.New(err_str)
	}
	secret_key := "a8bb5e292d234b5287c909cdba79d999"
	FBurl := fmt.Sprintf("GET https://graph.facebook.com/v2.8/oauth/access_token?client_id=%s&redirect_uri=%s&client_secret=%s&code=%s",
		"1788033858126569", "https://magicbackpack.ru/api/social/fb/getcode/", secret_key, code)

	resp, err := http.Get(FBurl)
	if err!=nil{
		return "", "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var answer FBTokenResponse
	json.Unmarshal(body, &answer)
	if answer.Access_token == "" {
		return "", "", errors.New("no token")
	}
	new_url := fmt.Sprintf("https://graph.facebook.com/v2.8/me?fields=id,name,email&access_token=%s",
		answer.Access_token,
	)
	new_resp, err := http.Get(new_url)
	if err!=nil{
		return "", "", err
	}
	defer new_resp.Body.Close()
	body, _ = ioutil.ReadAll(new_resp.Body)
	var user_info FBUserInfo
	json.Unmarshal(body, &user_info)

	var profile model.Profile
	if user_info.Email != ""{
		profile, err = model.GetProfileEmail(user_info.Email)
		if profile.Email == ""{
			profile.Email = user_info.Email
			profile.Name = user_info.Name
			profile.Password = randStringRunes(10)
			profile, err = model.CreateProfile(profile)
			if err != nil{
				return "", "", err
			}
		}
	} else {
		profile, err = model.GetProfileEmail(user_info.Id)
		if profile.Email == ""{
			profile.Email = user_info.Email
			profile.Name = user_info.Name
			profile.Password = randStringRunes(10)
			profile, err = model.CreateProfile(profile)
			if err != nil{
				return "", "", err
			}
		}
	}
	userToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  profile.ID,
		"type":"user",
	})
	bearToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  profile.ID,
		"type":"bear",
	})
	userTokenSigned, err := userToken.SignedString([]byte(hmacUserSecret))
	bearTokenSigned, err := bearToken.SignedString([]byte(hmacUserSecret))
	return userTokenSigned, bearTokenSigned, nil
}