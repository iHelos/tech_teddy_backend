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
	CreateStory_SpecifyID(Story{
		ID:14,
		Category:4,
		Name:"солнце и ветер",
		Price:0,
		Duration:"02:05",
		Description:"сказка о том, как солнце и ветер поспорили",
		AuthorID:1,
		Roled:true,
		DurationSplitted:Duration{
			Minutes:2,
			Seconds:5,
		},
		ImgUrls:UrlImage{
			Small:"https://storage.googleapis.com/hardteddy_images/small/14.jpg",
			Large:"https://storage.googleapis.com/hardteddy_images/large/14.jpg",
		},
		Parts:[]StoryPart{
			{
				Text:"Однажды сердитый северный Ветер и Солнце затеяли спор о том, кто из них сильнее. Долго они спорили и решили испробовать свою силу на одном путешественнике. Ветер сказал: -Я сейчас вмиг сорву с него плащ!",
				Part:"Часть 1",
				Audio:UrlAudio{
					Raw:"https://storage.googleapis.com/hardteddy_stories/14_1.raw",
				},
			},
			{
				Text:"И начал дуть. Он дул очень сильно и долго. Но человек только плотнее закутывался в свой плащ.",
				Part:"Часть 2",
				Audio:UrlAudio{
					Raw:"https://storage.googleapis.com/hardteddy_stories/14_2.raw",
				},
			},
			{
				Text:"Тогда Солнце начало пригревать путника. Он сначала опустил воротник, потом развязал пояс, а потом снял плащ и понёс его на руке.",
				Part:"Часть 3",
				Audio:UrlAudio{
					Raw:"https://storage.googleapis.com/hardteddy_stories/14_3.raw",
				},
			},
			{
				Text:"Солнце сказало Ветру: - Видишь: добром, лаской, можно добиться гораздо большего, чем насилием.",
				Part:"Часть 4",
				Audio:UrlAudio{
					Raw:"https://storage.googleapis.com/hardteddy_stories/14_4.raw",
				},
			},
		},
	})
}