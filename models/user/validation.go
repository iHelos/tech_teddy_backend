package user

import (
	"regexp"
	"github.com/asaskevich/govalidator"
)

const (
	regExpLogin string = "^[a-z0-9_-]{3,16}$"
	regExpPassword string = "^.{6,}$"
)

var (
	rxLogin = regexp.MustCompile(regExpLogin)
	rxPassword = regexp.MustCompile(regExpPassword)
)

func IsEmail(str string) bool {
	return govalidator.IsEmail(str)
}

func IsPassword(str string) bool{
	return rxPassword.MatchString(str)
}

func IsLogin(str string) bool{
	return rxLogin.MatchString(str)
}