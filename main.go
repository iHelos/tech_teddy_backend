package main

import (
	"github.com/kataras/iris"
	"github.com/tarantool/go-tarantool"
	"log"
	"time"
	"os"
	"github.com/iHelos/tech_teddy/sessionDB"
	"github.com/kataras/go-template/html"
	"github.com/iHelos/tech_teddy/filelogger"
	"github.com/iHelos/tech_teddy/teddyUsers"
	"github.com/iHelos/tech_teddy/teddyUsers/tarantool-user-storage"
)

var userstorage *teddyUsers.UserStorage

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
	userstorage = teddyUsers.New(iris.Config.Sessions.Cookie)
	userstorage.Engine = tarantool_user_storage.StorageConnection{client}

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

	saveapi := iris.Party("/saveapi/")
	saveapi.Use(filelogger.New("log.txt"))
	saveapi.Get("*randomName", func(ctx *iris.Context) {

	} )
	saveapi.Post("*randomName", func(ctx *iris.Context) {

	} )

	user := iris.Party("/user/")

	user.Get("/login", func(c *iris.Context) {
		c.Session().Set("name", "iris")
		c.Write("All ok session set to: %s", c.Session().GetString("name"))
	})

	user.Get("/registration", func(c *iris.Context) {
		name := c.Session().GetString("name")
		c.Write("The name on the /set was: %s", name)
	})

	user.Get("/logout", func(c *iris.Context) {
		c.SessionDestroy()
	})

	iris.Listen(":"+port)
}