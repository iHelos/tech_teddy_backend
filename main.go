package main

import (
	"github.com/kataras/iris"
	"github.com/tarantool/go-tarantool"
	"log"
	"time"
	"os"
	//sessionDB "github.com/iHelos/tech_teddy/models/session"
	"github.com/iHelos/tech_teddy/helper/filelogger"
	teddyUsers "github.com/iHelos/tech_teddy/models/user"
	"github.com/iHelos/tech_teddy/models/story"
	"github.com/iHelos/tech_teddy/models/user/tarantool-user-storage"
	"github.com/iHelos/tech_teddy/deploy-config"
	"github.com/iHelos/tech_teddy/helper/REST"
)

var userstorage *teddyUsers.UserStorage
var config *deploy_config.DeployConfiguration

func init() {
	config = deploy_config.GetConfiguration("./deploy.config")

	server := config.Database.Host
	opts := tarantool.Opts{
		Timeout:       500 * time.Millisecond,
		Reconnect:     1 * time.Second,
		MaxReconnects: 3,
		User:          config.Database.User,
		Pass:          config.Database.Password,
	}

	client, err := tarantool.Connect(server, opts)
	if err != nil {
		log.Fatalf("Failed to connect: %s", err.Error())
	}
	resp, err := client.Ping()
	log.Println(resp.Code)
	log.Println(resp.Data)
	log.Println(err)

	//sessionstorage := sessionDB.SessionConnection{client}

	userstorage = teddyUsers.New(iris.Config.Sessions.Cookie)
	userstorage.Engine = tarantool_user_storage.StorageConnection{client}
	//iris.UseSessionDB(sessionstorage)

	iris.Config.IsDevelopment = false
	iris.Config.Gzip = false
	iris.Config.Charset = "UTF-8"
	iris.Config.Sessions.DisableSubdomainPersistence = false
	iris.StaticServe("./static/web_files", "/static")

	//iris.UseTemplate(html.New(html.Config{
	//	Layout: "layout.html",
	//})).Directory("./templates", ".html")
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = config.Port
	}
	iris.Use(filelogger.New("logs/all.log"))
	iris.Get("/", func(ctx *iris.Context) {
		ctx.Render("index.html", nil)
	})

	iris.Get("/mock", func(ctx *iris.Context) {
		body := make([]map[string]string, 2)
		body[0] = map[string]string{
			"name":"iHelos",
			"email":"ihelos.ermakov@gmail.com",
		}
		body[1] = map[string]string{
			"name":"AnnJelly",
			"email":"annjellyiu5@gmail.com",
		}

		ctx.JSON(iris.StatusOK, REST.GetResponse(0, body))
	})

	iris.Get("/profile", userstorage.MustBeLogged, func(ctx *iris.Context) {

	})

	iris.Get("/story/:id", func(ctx *iris.Context) {
		id := ctx.Param("id")
		if id == "1" {
			ctx.SendFile("./static/audio/music.mp3", "music.mp3")
		} else {
			ctx.SendFile("./static/audio/story.mp3", "story.mp3")
		}

	})

	api := iris.Party("/api/")
	api.Get("/", func(ctx *iris.Context) {
		ctx.Redirect("http://docs.hardteddy.apiary.io")
	})
	saveapi := api.Party("/saveapi/")
	saveapi.Use(filelogger.New("logs/log.log"))
	saveapi.Get("*randomName", func(ctx *iris.Context) {

	})
	saveapi.Post("*randomName", func(ctx *iris.Context) {

	})

	api.Get("/allstories", func(ctx *iris.Context) {
		ctx.JSON(iris.StatusOK, REST.GetResponse(0, map[string]interface{}{
			"stories":story.GetAllStories(),
		}))
	})

	api.Get("/mystories", func(ctx *iris.Context) {
		ctx.JSON(iris.StatusOK, REST.GetResponse(0, map[string]interface{}{
			"stories":story.GetMyStories(),
		}))
	})

	// Пользовательские вьюхи
	apiuser := api.Party("/user/")
	apiuser.Use(filelogger.New("logs/userlog.log"))
	apiuser.Post("/login", func(ctx *iris.Context) {
		userToken, bearToken, err := userstorage.LoginUser(ctx)
		if err != nil {
			ctx.JSON(iris.StatusOK, REST.GetResponse(1, err.(*teddyUsers.UserError).Messages))
		} else {
			ctx.JSON(iris.StatusOK, REST.GetResponse(0, map[string]string{
				"userToken":userToken,
				"bearToken":bearToken,
			}))
		}
	})

	apiuser.Post("/register", func(ctx *iris.Context) {
		userToken, bearToken, err := userstorage.CreateUser(ctx)
		if err != nil {
			ctx.JSON(iris.StatusOK, REST.GetResponse(1, err.(*teddyUsers.UserError).Messages))
		}        else {
			ctx.JSON(iris.StatusOK, REST.GetResponse(0, map[string]string{
				"userToken":userToken,
				"bearToken":bearToken,
			}))
		}
	})("register")

	apiuser.Get("/sendall", func(ctx *iris.Context) {
		teddyUsers.SendAll()
		ctx.JSON(iris.StatusOK, REST.GetResponse(0, map[string]string{"":""}))
	})

	apiuser.Post("/logout", func(ctx *iris.Context) {
		ctx.SessionDestroy()
		ctx.JSON(iris.StatusOK, REST.GetResponse(0, map[string]string{"":""}))
	})

	iris.Listen(config.Host + ":" + port)
}
