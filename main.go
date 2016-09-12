package main

import (
	"github.com/kataras/iris"
	"github.com/tarantool/go-tarantool"
	"log"
	"time"
	"os"
	"github.com/iHelos/tech_teddy_backend/sessionDB"
	"github.com/kataras/go-template/html"
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
	iris.Config.IsDevelopment = true
	iris.Config.Gzip  = true
	iris.Config.Charset = "UTF-8"
	iris.UseTemplate(html.New(html.Config{
		Layout: "layout.html",
	})).Directory("./templates", ".html")
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	iris.Get("/", func(c *iris.Context) {
		c.Render("main.html", struct{Title string
					     Message []string}{"My Page title", []string{"message1","message2"}})

	})
	iris.Get("/set", func(c *iris.Context) {

		//set session values
		c.Session().Set("name", "iris")

		//test if setted here
		c.Write("All ok session setted to: %s", c.Session().GetString("name"))
	})

	iris.Get("/get", func(c *iris.Context) {
		// get a specific key, as string, if no found returns just an empty string
		//name := c.Session().GetString("name")

		c.Write("The name on the /set was: %s", c.Session().GetString("name"))
	})

	iris.Get("/delete", func(c *iris.Context) {
		// delete a specific key
		c.Session().Delete("name")
		//c.Session().ID()
	})

	iris.Get("/clear", func(c *iris.Context) {
		// removes all entries
		c.Session().Clear()
	})

	iris.Get("/destroy", func(c *iris.Context) {
		//destroy, removes the entire session and cookie
		c.SessionDestroy()
		c.Log("You have to refresh the page to completely remove the session (on browsers), so the name should NOT be empty NOW, is it?\n ame: %s\n\nAlso check your cookies in your browser's cookies, should be no field for localhost/127.0.0.1 (or what ever you use)", c.Session().GetString("name"))
		c.Write("You have to refresh the page to completely remove the session (on browsers), so the name should NOT be empty NOW, is it?\nName: %s\n\nAlso check your cookies in your browser's cookies, should be no field for localhost/127.0.0.1 (or what ever you use)", c.Session().GetString("name"))
	})



	iris.Listen(":"+port)
}