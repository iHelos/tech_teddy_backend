package view

import (
	"github.com/kataras/iris"
	"github.com/iHelos/tech_teddy/logic"
	"github.com/iHelos/tech_teddy/helper"
	"google.golang.org/api/option"
	"cloud.google.com/go/storage"
	"context"
	"log"
)

var google_client *storage.Client

func init(){
	ctx := context.Background()
	var err error
	google_client, err = storage.NewClient(
		ctx,
		option.WithServiceAccountFile("./gostorage.json"),
	)
	if err != nil {
		log.Print(err)
	}
}

func RenderPage (ctx *iris.Context) {
	ctx.Render("index.html", nil)
}

func UploadStory(ctx *iris.Context) {
	id := ctx.Param("id")
	logic.AddStoryFile(ctx, id, google_client)
}
func UploadSmallImage(ctx *iris.Context) {
	id := ctx.Param("id")
	logic.AddStorySmallImg(ctx, id, google_client)
}
func UploadLargeImage(ctx *iris.Context) {
	id := ctx.Param("id")
	logic.AddStoryLargeImg(ctx, id, google_client)
}

func ApiRedirect(ctx *iris.Context) {
	ctx.Redirect("http://docs.hardteddy.apiary.io")
}

func Login(ctx *iris.Context) {
	userToken, bearToken, err := logic.Login(ctx)
	if err != nil {
		ctx.JSON(iris.StatusOK, helper.GetResponse(1, err.(*helper.TeddyError).Messages))
	} else {
		ctx.JSON(iris.StatusOK, helper.GetResponse(0, map[string]string{
			"userToken":userToken,
			"bearToken":bearToken,
		}))
	}
}

func Register (ctx *iris.Context) {
	userToken, bearToken, err := logic.Register(ctx)
	if err != nil {
		ctx.JSON(iris.StatusOK, helper.GetResponse(1, err.(*helper.TeddyError).Messages))
	}        else {
		ctx.JSON(iris.StatusOK, helper.GetResponse(0, map[string]string{
			"userToken":userToken,
			"bearToken":bearToken,
		}))
	}
}

func MustBeLogged(ctx *iris.Context) {
	_,err := logic.ParseToken(ctx)
	if err != nil {
		ctx.JSON(iris.StatusOK, helper.GetResponse(
			1,
			map[string]int{
				"loginstatus":0,
			},
		))
	} else {
		ctx.Next()
	}
}

func GetUserStories(ctx *iris.Context) {
	stories, err := logic.GetMyStories(ctx)
	if (err != nil) {
		ctx.JSON(iris.StatusOK, helper.GetResponse(1, map[string]interface{}{
			"err":err.Error(),
		}))
	}else {
		ctx.JSON(iris.StatusOK, helper.GetResponse(0, map[string]interface{}{
			"stories":stories,
		}))
	}

}

func UserLikeStory(ctx *iris.Context) {
	err := logic.LikeStory(ctx)
	if (err != nil) {
		ctx.JSON(iris.StatusOK, helper.GetResponse(1, map[string]interface{}{
			"err":err.Error(),
		}))
	}else {
		ctx.JSON(iris.StatusOK, helper.GetResponse(0, map[string]interface{}{
		}))
	}
}
func GetStories(ctx *iris.Context) {
	stories, err := logic.GetStories(ctx)
	if (err != nil) {
		ctx.JSON(iris.StatusOK, helper.GetResponse(1, map[string]interface{}{
			"err":err.Error(),
		}))
	}else {
		ctx.JSON(iris.StatusOK, helper.GetResponse(0, map[string]interface{}{
			"stories":stories,
		}))
	}
}
func GetCategories(ctx *iris.Context) {
	categories, err := logic.GetCategories(ctx)
	if (err != nil) {
		ctx.JSON(iris.StatusOK, helper.GetResponse(1, map[string]interface{}{
			"err":err.Error(),
		}))
	}else {
		ctx.JSON(iris.StatusOK, helper.GetResponse(0, map[string]interface{}{
			"categories":categories,
		}))
	}
}
func CreateStory(ctx *iris.Context){

}
func Search(ctx *iris.Context) {
	stories, err := logic.FindStories(ctx)
	if (err != nil) {
		ctx.JSON(iris.StatusOK, helper.GetResponse(1, map[string]interface{}{
			"err":err.Error(),
		}))
	}else {
		ctx.JSON(iris.StatusOK, helper.GetResponse(0, map[string]interface{}{
			"stories":stories,
		}))
	}
}

func VKLoginPage(ctx *iris.Context){
	ctx.Redirect("https://oauth.vk.com/authorize?client_id=5806269&redirect_uri=https://magicbackpack.ru/api/social/vk/getcode&scope=4194304" )
}

type vkCode struct{
	Access_token string `json:"access_token"`
	User_id int	`json:"user_id"`
	Email string 	`json:"email"`
	Error string	`json:"error"`
}
func VKGetCode(ctx *iris.Context){
	userToken, bearToken, err := logic.VKGetCode(ctx)
	if err != nil {
		ctx.Redirect("https://magicbackpack.ru/api/social/error?err="+err.Error())
	} else {
		ctx.Redirect("https://magicbackpack.ru/api/social/success?userToken="+userToken+"&bearToken="+bearToken)
	}
}

func OKLoginPage(ctx *iris.Context){
	ctx.Redirect("https://connect.ok.ru/oauth/authorize?client_id=1249370880&scope=GET_EMAIL&response_type=code&redirect_uri=https://magicbackpack.ru/api/social/ok/getcode&layout=m" )
}

func OKGetCode(ctx *iris.Context){
	userToken, bearToken, err := logic.OKGetCode(ctx)
	if err != nil {
		ctx.Redirect("https://magicbackpack.ru/api/social/error?err="+err.Error())
	} else {
		ctx.Redirect("https://magicbackpack.ru/api/social/success?userToken="+userToken+"&bearToken="+bearToken)
	}
}


func FBLoginPage(ctx *iris.Context){
	ctx.Redirect("https://www.facebook.com/v2.8/dialog/oauth?client_id=1788033858126569&redirect_uri=https://magicbackpack.ru/api/social/fb/getcode/&scope=email" )
}

func FBGetCode(ctx *iris.Context){
	userToken, bearToken, err := logic.FBGetCode(ctx)
	if err != nil {
		ctx.Redirect("https://magicbackpack.ru/api/social/error?err="+err.Error())
	} else {
		ctx.Redirect("https://magicbackpack.ru/api/social/success?userToken="+userToken+"&bearToken="+bearToken)
	}
}