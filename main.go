package main

import (
	"github.com/kataras/iris"
	"github.com/tarantool/go-tarantool"
	"log"
	"time"
	"os"
	"github.com/iHelos/tech_teddy/sessionDB"
	"github.com/kataras/go-template/html"
	"github.com/iris-contrib/middleware/logger"
)

type sessionConnection struct{
	*tarantool.Connection
}

func init()  {
	server := "77.244.214.4:3301"
	opts := tarantool.Opts{
		Timeout:       500 * time.Millisecond,
		Reconnect:     1 * time.Second,
		MaxReconnects: 3,
		User:          "goClient",
		Pass:          "TeddyTarantoolS1cret",
	}

	client, err := tarantool.Connect(server, opts)
	if err != nil {
		log.Fatalf("Failed to connect: %s", err.Error())
	}
	resp, err := client.Ping()
	log.Println(resp.Code)
	log.Println(resp.Data)
	log.Println(err)

	sessionstorage := sessionDB.SessionConnection{client}

	iris.UseSessionDB(sessionstorage)
	iris.Config.IsDevelopment = false
	iris.Config.Gzip  = false
	iris.Config.Charset = "UTF-8"
	iris.Config.Sessions.DisableSubdomainPersistence = false
	iris.StaticServe("./static")

	iris.UseTemplate(html.New(html.Config{
		Layout: "layout.html",
	})).Directory("./templates", ".html")
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	saveapi := iris.Party("/saveapi")

	mylogs := logger.New()

	saveapi.Get("/*randomName", func(ctx *iris.Context) {
		mylogs.Serve(ctx)
	} )

	saveapi.Post("/*randomName", func(ctx *iris.Context) {
		mylogs.Serve(ctx)
	} )

	iris.Listen(":"+port)
}