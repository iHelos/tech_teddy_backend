package model

import (
	"github.com/tarantool/go-tarantool"
	"gopkg.in/vmihailenco/msgpack.v2"
	"log"
	"reflect"
)

var client *tarantool.Connection = nil

func init() {
	msgpack.Register(reflect.TypeOf(Profile{}), encodeProfile, decodeProfile)
}

func InitDB(server string, opts tarantool.Opts){
	var err error
	client, err = tarantool.Connect(server, opts)
	if err != nil {
		log.Fatalf("Failed to connect: %s", err.Error())
	}
	asd := Profile{ID:1, Name:"wow", Email:"het", Likes:[]int{1,2,3}}
	kek, err := updateProfile(asd)
	log.Print(kek)
	kek,err = getProfile(1)
	log.Print(err)
}