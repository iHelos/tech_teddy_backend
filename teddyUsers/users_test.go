package teddyUsers

import (
	"testing"
)

type errorStruct struct {
	error
}


func TestLoginInvalid(t *testing.T) {
	var user = NewUser{
		User: User{Login: "", Email:"ihelos.ermakov@gmail.com"},
		Password1:"password",
		Password2:"password",
	}
	err := user.check()
	t.Log(err)
	if err == nil{
		t.Errorf("User Login error checking working bad")
	}
}

func TestEmailInvalid(t *testing.T) {
	var user = NewUser{
		User: User{Login: "ihelos", Email:"ihelos.ermakovgmail.com"},
		Password1:"password",
		Password2:"password",
	}
	err := user.check()
	t.Log(err)
	if err == nil{
		t.Errorf("User Email error checking working bad")
	}
}

func TestPasswordInvalid(t *testing.T) {
	var user = NewUser{
		User: User{Login: "ihelos", Email:"ihelos.ermakov@gmail.com"},
		Password1:"",
		Password2:"password",
	}
	err := user.check()
	t.Log(err)
	if err == nil{
		t.Errorf("User Password error checking working bad")
	}
}

func TestSecondPasswordInvalid(t *testing.T) {
	var user = NewUser{
		User: User{Login: "ihelos", Email:"ihelos.ermakov@gmail.com"},
		Password1:"password",
		Password2:"password2",
	}
	err := user.check()
	t.Log(err)
	if err == nil{
		t.Errorf("User Match passwords error checking working bad")
	}
}

func TestAllErrorsAtOnce(t *testing.T) {
	var user = NewUser{
		User: User{Login: "", Email:"ihelos.ermakovgmail.com"},
		Password1:"",
		Password2:"password",
	}
	err := user.check()
	t.Log(err)
	if err == nil{
		t.Errorf("Right user")
	}
}

func TestAllRight(t *testing.T) {
	var user = NewUser{
		User: User{Login: "ihelos", Email:"ihelos.ermakov@gmail.com"},
		Password1:"password",
		Password2:"password",
	}
	err := user.check()
	t.Log(err)
	if err != nil{
		t.Errorf("User checking working bad")
	}
}