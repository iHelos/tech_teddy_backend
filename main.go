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
	"github.com/iHelos/tech_teddy/models/user/tarantool-user-storage"
	teddyStory"github.com/iHelos/tech_teddy/models/story"
	"github.com/iHelos/tech_teddy/models/story/tarantool-story-storage"
	"github.com/iHelos/tech_teddy/deploy-config"
	"github.com/iHelos/tech_teddy/helper/REST"
	"github.com/iris-contrib/middleware/cors"
	"github.com/iris-contrib/middleware/recovery"
	"github.com/iHelos/tech_teddy/views/store"
)

var userstorage *teddyUsers.UserStorage
var storystorage teddyStory.StoryStorageEngine
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
	//TODO remake userStorage
	userstorage = teddyUsers.New(iris.Config.Sessions.Cookie)
	userstorage.Engine = tarantool_user_storage.StorageConnection{client}
	storystorage = tarantool_story_storage.StorageConnection{client}
	//iris.UseSessionDB(sessionstorage)

	iris.Use(recovery.New())
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

	cors_config := cors.Options{
		AllowedOrigins:[]string{"*"},
		AllowedMethods:[]string{"GET", "POST", "OPTIONS", ""},
		AllowCredentials:true,
		MaxAge:5,
		Debug:false,
	}

	cors_obj := cors.New(cors_config)

	iris.Use(cors_obj)
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
		switch id {
		case "1": ctx.SendFile("./static/audio/music.mp3", "music.mp3");
		case "2": ctx.SendFile("./static/audio/music.mp3", "music.mp3");
		case "3": ctx.SendFile("./static/audio/k_r.raw", "3.raw");
		case "4": ctx.SendFile("./static/audio/m_i_m.raw", "4.raw");
		case "5": ctx.SendFile("./static/audio/tale.raw", "5.raw");
		default:  ctx.SendFile("./static/audio/k_r.raw", "3.raw");
		}
	})

	api := iris.Party("/api/")
	api.Get("/", func(ctx *iris.Context) {
		ctx.Redirect("http://docs.hardteddy.apiary.io")
	})


	// Пользовательские вьюхи
	apiuser := api.Party("/user/")
	apiuser.Use(filelogger.New("logs/userlog.log"))
	apiuser.Any("/login", func(ctx *iris.Context) {
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

	apiuser.Get("/mystories", userstorage.MustBeLogged, func(ctx *iris.Context) {
		stories, err := store.GetMyStories(ctx, &storystorage)
		if (err != nil) {
			ctx.JSON(iris.StatusOK, REST.GetResponse(1, map[string]interface{}{
				"err":err.Error(),
			}))
		}else {
			ctx.JSON(iris.StatusOK, REST.GetResponse(0, map[string]interface{}{
				"stories":stories,
			}))
		}

	})

	apiuser.Any("/register", func(ctx *iris.Context) {
		userToken, bearToken, err := userstorage.CreateUser(ctx)
		if err != nil {
			ctx.JSON(iris.StatusOK, REST.GetResponse(1, err.(*teddyUsers.UserError).Messages))
		}        else {
			ctx.JSON(iris.StatusOK, REST.GetResponse(0, map[string]string{
				"userToken":userToken,
				"bearToken":bearToken,
			}))
		}
	})

	apiuser.Get("/sendall", func(ctx *iris.Context) {
		teddyUsers.SendAll()
		ctx.JSON(iris.StatusOK, REST.GetResponse(0, map[string]string{"":""}))
	})

	apiuser.Any("/logout", func(ctx *iris.Context) {
		ctx.SessionDestroy()
		ctx.JSON(iris.StatusOK, REST.GetResponse(0, map[string]string{"":""}))
	})

	apistore := api.Party("/store/")
	apistore.Get("/stories", func(ctx *iris.Context) {
		stories, err := store.GetAllStories(ctx, &storystorage)
		if (err != nil) {
			ctx.JSON(iris.StatusOK, REST.GetResponse(1, map[string]interface{}{
				"err":err.Error(),
			}))
		}else {
			ctx.JSON(iris.StatusOK, REST.GetResponse(0, map[string]interface{}{
				"stories":stories,
			}))
		}
	})

	iris.Listen(config.Host + ":" + port)
}
