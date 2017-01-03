package main

import (
	"github.com/kataras/iris"
	"github.com/tarantool/go-tarantool"
	"log"
	"time"
	"os"
	//sessionDB "github.com/iHelos/tech_teddy/models/session"
	"github.com/iHelos/tech_teddy/deploy-config"
	"github.com/iris-contrib/middleware/cors"
	//"github.com/iris-contrib/middleware/recovery"
	"cloud.google.com/go/storage"
	"context"
	"google.golang.org/api/option"
	"fmt"
	"io/ioutil"
	"github.com/iHelos/tech_teddy/model"
	"github.com/iHelos/tech_teddy/helper"
	"github.com/iHelos/tech_teddy/view"
)

var config *deploy_config.DeployConfiguration
var google_client *storage.Client
var cors_obj *cors.Cors

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

	model.InitDB(server,opts)

	//iris.Use(recovery.New())
	iris.Config.IsDevelopment = false
	iris.Config.Gzip = false
	iris.Config.Charset = "UTF-8"
	iris.Config.Sessions.DisableSubdomainPersistence = false
	iris.StaticServe("./static/web_files", "/static")

	ctx := context.Background()
	google_client, err := storage.NewClient(
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

	iris.Use(helper.New("logs/all.log"))

	iris.Get("/", view.RenderPage)

	iris.Any("/upload/story/:id", view.UploadStory)
	iris.Any("/upload/smallimg/:id", view.UploadSmallImage)
	iris.Any("/upload/largeimg/:id", view.UploadLargeImage)

	//iris.Get("/story/:id", func(ctx *iris.Context) {
	//	id := ctx.Param("id")
	//	ctx.SendFile("./static/audio/" + id + ".raw", id + ".raw")
	//})

	api := iris.Party("/api/")
	api.Get("/", view.ApiRedirect)

	// Пользовательские вьюхи
	apiuser := api.Party("/user/")
	apiuser.Use(helper.New("logs/userlog.log"))
	apiuser.Post("/login", view.Login)
	apiuser.Post("/signup", view.Register)
	apiuser.Get("/mystories", view.MustBeLogged, view.GetUserStories)

	apistore := api.Party("/store/")
	//apistore.Any("/story/add", func(ctx *iris.Context) {
	//	story_obj, err := store.AddStory(ctx, &storystorage)
	//	if (err != nil) {
	//		ctx.JSON(iris.StatusOK, REST.GetResponse(1, map[string]interface{}{
	//			"err":err.Error(),
	//		}))
	//	}else {
	//		ctx.JSON(iris.StatusOK, REST.GetResponse(0, map[string]interface{}{
	//			"story":story_obj,
	//		}))
	//	}
	//})
	apistore.Get("/story/", view.GetStories)
	apistore.Get("/categories/", view.GetCategories)
	apistore.Any("/buy", view.UserLikeStory)
	apistore.Get("/search/", view.Search)

	iris.Set(iris.OptionMaxRequestBodySize(64 << 20))
	iris.Listen(config.Host + ":" + port)
	//iris.ListenLETSENCRYPT(config.Host + ":" + port)
}
