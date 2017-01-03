package model

import (
	"fmt"
	"reflect"
	"gopkg.in/vmihailenco/msgpack.v2"
	"github.com/tarantool/go-tarantool"
	"errors"
	"golang.org/x/crypto/bcrypt"
	. "github.com/iHelos/tech_teddy/helper"
	"regexp"
	"github.com/asaskevich/govalidator"
	"github.com/labstack/gommon/log"
)

type Profile struct {
	ID       int
	Email    string
	Name     string
	Password string
	Likes    []int
	Bears    []int
}

func encodeProfile(e *msgpack.Encoder, v reflect.Value) error {
	m := v.Interface().(Profile)
	if m.ID > 0 {
		if err := e.EncodeArrayLen(6); err != nil {
			return err
		}
		if err := e.EncodeInt(m.ID); err != nil {
			return err
		}
	} else{
		if err := e.EncodeArrayLen(5); err != nil {
			return err
		}
	}
	if err := e.EncodeString(m.Email); err != nil {
		return err
	}
	if err := e.EncodeString(m.Name); err != nil {
		return err
	}
	if err := e.EncodeString(m.Password); err != nil {
		return err
	}
	if err := e.EncodeArrayLen(len(m.Likes)); err != nil {
		return err
	}
	for _, l := range m.Likes {
		e.EncodeInt(l)
	}
	if err := e.EncodeArrayLen(len(m.Bears)); err != nil {
		return err
	}
	for _, b := range m.Likes {
		e.EncodeInt(b)
	}
	return nil
}

func decodeProfile(d *msgpack.Decoder, v reflect.Value) error {
	var err error
	var l int
	m := v.Addr().Interface().(*Profile)
	if l, err = d.DecodeArrayLen(); err != nil {
		return err
	}
	if l != 6 {
		return fmt.Errorf("array len doesn't match: %d", l)
	}
	if m.ID, err = d.DecodeInt(); err != nil {
		return err
	}
	if m.Email, err = d.DecodeString(); err != nil {
		return err
	}
	if m.Name, err = d.DecodeString(); err != nil {
		return err
	}
	if m.Password, err = d.DecodeString(); err != nil {
		return err
	}
	if l, err = d.DecodeArrayLen(); err != nil {
		return err
	}
	m.Likes = make([]int, l)
	for i := 0; i < l; i++ {
		m.Likes[i],_ = d.DecodeInt()
	}
	if l, err = d.DecodeArrayLen(); err != nil {
		return err
	}
	m.Bears = make([]int, l)
	for i := 0; i < l; i++ {
		m.Bears[i],_ = d.DecodeInt()
	}
	return nil
}

func CreateProfile(new_profile Profile) (created_profile Profile, err error){
	var profiles []Profile
	err = client.Call17Typed("box.space.user:auto_increment", []interface{}{new_profile}, &profiles)
	if err!=nil {
		return Profile{}, err
	}
	return profiles[0], err
}

func UpdateProfile(new_profile Profile) (updated_profile Profile, err error){
	var profiles []Profile
	err = client.ReplaceTyped("user", new_profile, &profiles)
	return profiles[0], err
}

func GetProfile(id int) (profile Profile, err error){
	var profiles []Profile
	err = client.SelectTyped("user", "primary", 0,1, tarantool.IterEq, []interface{}{uint(id)}, &profiles)
	if len(profiles)>0 {
		return profiles[0], err
	}
	return profile, errors.New("not exists")
}

func CheckPassword(profile Profile) (int, error){
	var dbProfile []Profile
	err := client.SelectTyped("user", "email", 0, 1, tarantool.IterEq, []interface{}{profile.Email}, &dbProfile)
	if err != nil {
		return 0, err
	}
	log.Print(dbProfile)
	err = bcrypt.CompareHashAndPassword([]byte(dbProfile[0].Password), []byte(profile.Password))
	if err != nil{
		return 0, err
	}
	return dbProfile[0].ID, nil
}

const (
	regExpLogin string = "^[a-z0-9_-]{3,16}$"
	regExpPassword string = "^.{6,}$"
)

var (
	rxLogin = regexp.MustCompile(regExpLogin)
	rxPassword = regexp.MustCompile(regExpPassword)
)

type NewProfile struct  {
	Login string `json:"name" form:"name"`
	Email string `json:"email" form:"email"`
	Password1 string `json:"password1" form:"password1"`
	Password2 string `json:"password2" form:"password2"`
}

func IsEmail(str string) bool {
	return govalidator.IsEmail(str)
}

func IsPassword(str string) bool{
	return rxPassword.MatchString(str)
}

func IsLogin(str string) bool{
	return rxLogin.MatchString(str)
}

func (user *NewProfile) Validate() error{
	var err = NewError()
	if ok := IsLogin(user.Login); !ok {
		err.Append("login", 0)
	}
	if ok := IsEmail(user.Email); !ok {
		err.Append("email", 0)
	}
	if ok := IsPassword(user.Password1); !ok {
		err.Append("password", 0)
	}
	if user.Password1 != user.Password2{
		err.Append("password", 1)
	}
	if len(err.Messages) != 0 {return err}
	return nil
}

func AddStory(id int, sid int) (error){
	_, err := client.Call("AddStory", []interface{}{id, sid})
	return err
}