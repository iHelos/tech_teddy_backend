package main_test

import (
	"log"
	"io/ioutil"
	"testing"
	"net/http"
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

