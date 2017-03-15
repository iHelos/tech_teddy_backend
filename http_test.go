package main_test

import (
	"log"
	"io/ioutil"
	"testing"
	"net/http"
	."github.com/iHelos/tech_teddy/model"
	"fmt"
)

func BenchmarkHttpParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := http.Get("http://localhost:8080/cookie/set")
			if err != nil {
				log.Fatal(err)
			}
			_, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			//b.Log(string(body))
		}
	})
}

func ExampleReverse() {
	story := Story{
		ID:1,
		Name:"баба яга",
		Category:1,
		Price:0,
		Duration:"05:55",
		Description:"Жили-были муж с женой, и была у них дочка. Заболела жена и умерла. Погоревал-погоревал мужик да и женился на другой.Невзлюбила злая баба девочку, била ее, ругала, только и думала, как бы совсем извести, погубить.",
		AuthorID:1,
		Roled:false,
		DurationSplitted:Duration{
			Minutes:5,
			Seconds:55,
		},
		ImgUrls:UrlImage{
			Small:"http://storage.googleapis.com/hardteddy_images/small/1.jpg",
			Large:"http://storage.googleapis.com/hardteddy_images/large/1.jpg",
		},
		Parts:[]StoryPart{
			{
				Part:"1",
				Text:"",
				Audio:UrlAudio{
					Raw:"http://storage.googleapis.com/hardteddy_stories/1.raw",
					Original:"http://storage.googleapis.com/hardteddy_stories/mp3/1.mp3",
				},
				Size:2843501,
				Role:"",
			},
		},
		Roles:[]string{},
	}
	UpdateStory(story)
	fmt.Println(story.Name)
	// Output:
	// баба яга
}

func ExampleStory() {
	story := Story{
		ID:15,
		Name:"солнце и ветер часть 2. возвращение путника",
		Category:4,
		Price:0,
		Duration:"02:05",
		Description:"сказка о том, как солнце и ветер поспорили",
		AuthorID:1,
		Roled:true,
		Roles:[]string{"ветер", "солнце", "рассказчик"},
		DurationSplitted:Duration{
			Minutes:2,
			Seconds:5,
		},
		ImgUrls:UrlImage{
			Small:"http://storage.googleapis.com/hardteddy_images/small/14.jpg",
			Large:"http://storage.googleapis.com/hardteddy_images/large/14.jpg",
		},
		Parts:[]StoryPart{
			{
				Part:"Часть 1",
				Text:"Однажды сердитый северный Ветер и Солнце затеяли спор о том, кто из них сильнее. Долго они спорили и решили испробовать свою силу на одном путешественнике. Ветер сказал: -Я сейчас вмиг сорву с него плащ!",
				Audio:UrlAudio{
					Raw:"http://storage.googleapis.com/hardteddy_stories/14_1.raw",
					Original:"http://storage.googleapis.com/hardteddy_stories/mp3/14_1.mp3",
				},
				Role:"рассказчик",
				Size:148570,
			},
			{
				Part:"Часть 2",
				Text:"И начал дуть. Он дул очень сильно и долго. Но человек только плотнее закутывался в свой плащ.",
				Audio:UrlAudio{
					Raw:"http://storage.googleapis.com/hardteddy_stories/14_2.raw",
					Original:"http://storage.googleapis.com/hardteddy_stories/mp3/14_2.mp3",
				},
				Role:"ветер",
				Size:72860,
			},
			{
				Part:"Часть 3",
				Text:"Тогда Солнце начало пригревать путника. Он сначала опустил воротник, потом развязал пояс, а потом снял плащ и понёс его на руке.",
				Audio:UrlAudio{
					Raw:"http://storage.googleapis.com/hardteddy_stories/14_3.raw",
					Original:"http://storage.googleapis.com/hardteddy_stories/mp3/14_3.mp3",
				},
				Role:"рассказчик",
				Size:115510,
			},
			{
				Part:"Часть 4",
				Text:"Солнце сказало Ветру: - Видишь: добром, лаской, можно добиться гораздо большего, чем насилием.",
				Audio:UrlAudio{
					Raw:"http://storage.googleapis.com/hardteddy_stories/14_4.raw",
					Original:"http://storage.googleapis.com/hardteddy_stories/mp3/14_4.mp3",
				},
				Role:"солнце",
				Size:85510,
			},
		},
	}
	UpdateStory(story)
	fmt.Println(story.Name)
	// Output:
	// баба яга
}