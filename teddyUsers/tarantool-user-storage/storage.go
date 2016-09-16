package tarantool_user_storage

import (
	"github.com/iHelos/tech_teddy/teddyUsers"
	"github.com/tarantool/go-tarantool"
	"log"

	"golang.org/x/crypto/bcrypt"

	"errors"
)

type StorageConnection struct {
	*tarantool.Connection
}

func (con StorageConnection) Create(name, email, password string) (error) {
	hashpassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	_, err := con.Call("createProfile", []interface{}{name, email, hashpassword})
	return err
}

func (con StorageConnection) Load(sid string) (teddyUsers.NewUser, error){
	data, err := con.Call("getProfile", []interface{}{sid})
	if err != nil {
		return teddyUsers.NewUser{}, err
	}
	dataslice := data.Data[0].([]interface{})
	if len(dataslice) < 4{
		return teddyUsers.NewUser{}, errors.New("invalid user")
	} else{}

	return teddyUsers.NewUser{}, nil
}

func (con StorageConnection) CheckLogin(login, password string) (error){
	resp, err := con.Select("profile", "primary", 0, 1, tarantool.IterEq, []interface{}{login})
	if err != nil {
		return err
	}
	if len(resp.Data) != 1{
		return errors.New("no such user")
	}
	dataslice := resp.Data[0].([]interface{})
	if len(dataslice) < 3{
		return errors.New("user invalid")
	}
	if hashedPass, ok := dataslice[2].(string); ok{
		log.Print(hashedPass)
		err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password))
		if err != nil{
			return err
		}
		return nil
	}
	return errors.New("user invalid")
}

func (con StorageConnection) CheckIsLogged(sid string) (error){
	_, err := con.Call("isLogined", []interface{}{sid})
	return err
}
