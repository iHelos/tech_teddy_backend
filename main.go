package main

import (
	"github.com/kataras/iris"
	"github.com/tarantool/go-tarantool"
	"log"
	"time"
	"os"
	"github.com/iHelos/tech_teddy/sessionDB"
//	"github.com/kataras/go-template/html"
	"github.com/iHelos/tech_teddy/filelogger"
	"github.com/iHelos/tech_teddy/teddyUsers"
	"github.com/iHelos/tech_teddy/teddyUsers/tarantool-user-storage"
	"github.com/asaskevich/govalidator"
	"strings"
)

var userstorage *teddyUsers.UserStorage

func init() {
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
	iris.Config.IsDevelopment = true
	iris.Config.Gzip = false
	iris.Config.Charset = "UTF-8"
	iris.Config.Sessions.DisableSubdomainPersistence = false
	iris.StaticServe("./static")

	//iris.UseTemplate(html.New(html.Config{
	//	Layout: "layout.html",
	//})).Directory("./templates", ".html")
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	iris.Get("/", func(ctx *iris.Context) {
		ctx.Render("index.html", nil)
	})

	iris.Get("/profile", userstorage.MustBeLogged, func(ctx *iris.Context){
		err := userstorage.LoginUser(ctx)
		log.Print(err)
	})

	api := iris.Party("/api/")
	api.Get("/", func(ctx *iris.Context) {
		ctx.Redirect("http://docs.hardteddy.apiary.io")
	})
	saveapi := api.Party("/saveapi/")
	saveapi.Use(filelogger.New("log.log"))
	saveapi.Get("*randomName", func(ctx *iris.Context) {

	})
	saveapi.Post("*randomName", func(ctx *iris.Context) {

	})

	apiuser := api.Party("/user/")
	apiuser.Use(filelogger.New("userlog.log"))
	apiuser.Post("/login", func(ctx *iris.Context) {
		err := userstorage.LoginUser(ctx)
		if err != nil{
			ctx.JSON(iris.StatusOK, map[string]string{"error":err.Error()})
		} else {
			ctx.JSON(iris.StatusOK, map[string]string{"sessionid":ctx.Session().ID()})
		}
	})

	apiuser.Post("/register", func(ctx *iris.Context) {
		err := userstorage.CreateUser(ctx)
		if errors, ok := err.(govalidator.Errors); ok {
			errs := make(map[string]string)
			log.Print(errors)
			for _, msg := range errors {
				values := strings.Split(msg.Error(), ":")
				errs[values[0]] = values[1]
				log.Print(msg)
			}
			ctx.JSON(iris.StatusOK, errs)
		} else if err != nil {
			ctx.JSON(iris.StatusOK, map[string]string{"error":err.Error()})
		}        else {
			ctx.JSON(iris.StatusOK, map[string]string{"sessionid":ctx.Session().ID()} )
		}
	})("register")

	apiuser.Post("/logout", func(ctx *iris.Context) {
		ctx.SessionDestroy()
	})

	iris.Listen(":" + port)
}