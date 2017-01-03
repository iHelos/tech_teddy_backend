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
	msgpack.Register(reflect.TypeOf(Story{}), encodeStory, decodeStory)
	msgpack.Register(reflect.TypeOf(StoryPart{}), encodeStoryPart, decodeStoryPart)
	msgpack.Register(reflect.TypeOf(Duration{}), encodeDuration, decodeDuration)
	msgpack.Register(reflect.TypeOf(UrlImage{}), encodeUrlImage, decodeUrlImage)
	msgpack.Register(reflect.TypeOf(UrlAudio{}), encodeUrlAudio, decodeUrlAudio)
}

func InitDB(server string, opts tarantool.Opts){
	var err error
	client, err = tarantool.Connect(server, opts)
	if err != nil {
		log.Fatalf("Failed to connect: %s", err.Error())
	}
}