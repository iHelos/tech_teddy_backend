package main_test

import (
	"log"
	"io/ioutil"
	"testing"
	"net/http"
	."github.com/iHelos/tech_teddy/model"
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
		Duration:"5:55",
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
			},
		},
	}
	UpdateStory(story)
}
