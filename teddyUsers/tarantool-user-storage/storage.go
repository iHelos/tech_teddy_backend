package tarantool_user_storage

import (
	"github.com/iHelos/tech_teddy/teddyUsers"
	"github.com/tarantool/go-tarantool"
	"log"
)

type StorageConnection struct {
	*tarantool.Connection
}

func (con StorageConnection) Create(name, email, password string) (error) {
	ans, err := con.Call("createProfile", []interface{}{name, email, password})
	log.Print(ans)
	log.Print(err)
	return err
}

func (con StorageConnection) Load(string) (teddyUsers.NewUser, error){
	return teddyUsers.NewUser{}, nil
}
