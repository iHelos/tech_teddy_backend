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
	"cloud.google.com/go/storage"
	"context"
	"google.golang.org/api/option"
	"fmt"
	"io/ioutil"
)

var userstorage *teddyUsers.UserStorage
var storystorage teddyStory.StoryStorageEngine
var config *deploy_config.DeployConfiguration
var google_client *storage.Client

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
	ctx := context.Background()
	google_client, err = storage.NewClient(
		ctx,
		option.WithServiceAccountFile("./gostorage.json"),
	)
	if err != nil {
		log.Print(err)
	}
	_ = google_client
	bkt := google_client.Bucket("hardteddy_stories")
	attrs, err := bkt.Attrs(ctx)
	if err != nil {
		log.Print(err)
	}
	fmt.Printf("bucket %s, created at %s, is located in %s with storage class %s\n",
		attrs.Name, attrs.Created, attrs.Location, attrs.StorageClass)
	storage_opts := storage.SignedURLOptions{}
	storage_opts.PrivateKey, err = ioutil.ReadFile("./gostorage.pem")
	if err != nil{
		log.Print(err)
	}
	storage_opts.Expires = time.Now().Add(time.Minute)
	storage_opts.Method = "GET"
	storage_opts.GoogleAccessID = "116466809002114199830"
	url, err := storage.SignedURL("hardteddy_stories", "report.pdf", &storage_opts)
	if err != nil {
		log.Print(err)
	}
	fmt.Println(url)
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

	iris.Any("/upload/story/:id", func(ctx *iris.Context) {
		id := ctx.Param("id")
		store.AddStoryFile(ctx, id, google_client)
	})
	iris.Any("/upload/smallimg/:id", func(ctx *iris.Context) {
		id := ctx.Param("id")
		store.AddStorySmallImg(ctx, id, google_client)
	})
	iris.Any("/upload/largeimg/:id", func(ctx *iris.Context) {
		id := ctx.Param("id")
		store.AddStoryLargeImg(ctx, id, google_client)
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
		ctx.SendFile("./static/audio/" + id + ".raw", id + ".raw")
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

	apiuser.Get("/VKlogin", func(ctx *iris.Context) {
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
	apistore.Post("/story/add", func(ctx *iris.Context) {
		story_obj, err := store.AddStory(ctx, &storystorage)
		if (err != nil) {
			ctx.JSON(iris.StatusOK, REST.GetResponse(1, map[string]interface{}{
				"err":err.Error(),
			}))
		}else {
			ctx.JSON(iris.StatusOK, REST.GetResponse(0, map[string]interface{}{
				"story":story_obj,
			}))
		}
	})

	apistore.Get("/story/", func(ctx *iris.Context) {
		stories, err := store.GetStories(ctx, &storystorage)
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

	apistore.Get("/categories/", func(ctx *iris.Context) {
		categories, err := store.GetCategories(ctx, &storystorage)
		if (err != nil) {
			ctx.JSON(iris.StatusOK, REST.GetResponse(1, map[string]interface{}{
				"err":err.Error(),
			}))
		}else {
			ctx.JSON(iris.StatusOK, REST.GetResponse(0, map[string]interface{}{
				"categories":categories,
			}))
		}
	})

	apistore.Any("/buy", func(ctx *iris.Context) {
		_,err := userstorage.Buy(ctx);
		if err != nil {
			ctx.JSON(iris.StatusOK, REST.GetResponse(1, err))
		}        else {
			ctx.JSON(iris.StatusOK, REST.GetResponse(0, map[string]string{
				"transaction":"completed",
			}))
		}
	})

	apistore.Get("/search/", func(ctx *iris.Context) {
		stories, err := store.FindStories(ctx, &storystorage)
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

	iris.Set(iris.OptionMaxRequestBodySize(64 << 20))
	iris.Listen(config.Host + ":" + port)
}
