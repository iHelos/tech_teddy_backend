package store

import (
	"github.com/kataras/iris"
	"github.com/iHelos/tech_teddy/models/story"
	"github.com/iHelos/tech_teddy/models/user"
	"github.com/labstack/gommon/log"
	"strconv"
	"strings"
	"os"
	"fmt"
	"io"
	"mime/multipart"
//	"github.com/bobertlo/go-mpg123/mpg123"
	"os/exec"
	"encoding/json"
	"cloud.google.com/go/storage"
	"context"
)

func BuyStory(ctx *iris.Context, storage *story.StoryStorageEngine) ([]story.Story, error) {
	var stories = []story.Story{}
	login, err := user.GetLogin(ctx)
	if (err != nil){
		return stories, err
	}
	stories, err = (*storage).GetMyStories(login)
	return stories, err
}

func GetMyStories(ctx *iris.Context, storage *story.StoryStorageEngine) ([]story.Story, error) {
	var stories = []story.Story{}
	login, err := user.GetLogin(ctx)
	if (err != nil){
		log.Print(login, err)
		return stories, err
	}
	stories, err = (*storage).GetMyStories(login)
	return stories, err
}

func FindStories(ctx *iris.Context, storage *story.StoryStorageEngine) ([]story.Story, error) {
	var stories = []story.Story{}
	keyword := ctx.FormValueString("keyword")
	if len(keyword) < 3{
		return stories, nil
	}
	keyword = strings.ToLower(keyword)
	log.Print(keyword)
	stories, err := (*storage).Search(keyword)
	return stories, err
}

type StoriesParams struct{
	Cat int `form:"cat"`
	Page int `form:"page"`
	Order string `form:"order"`
	Order_Type string `form:"ordtype"`
}

func GetStories(ctx *iris.Context, storage *story.StoryStorageEngine) ([]story.Story, error){
	getStoriesParams := StoriesParams{}
	getStoriesParams.Cat, _ = strconv.Atoi(ctx.FormValueString("cat"))
	getStoriesParams.Page, _ = strconv.Atoi(ctx.FormValueString("page"))
	getStoriesParams.Order = ctx.FormValueString("order")
	getStoriesParams.Order_Type = ctx.FormValueString("ordtype")
	var stories = []story.Story{}
	var err error
	if (getStoriesParams.Cat == 0){
		stories, err = (*storage).GetAll(getStoriesParams.Order, getStoriesParams.Order_Type, getStoriesParams.Page)
	}else {
		stories, err = (*storage).GetAllByCategory(getStoriesParams.Order, getStoriesParams.Order_Type, getStoriesParams.Page, getStoriesParams.Cat)
	}
	return stories,err
}

type Category struct {
	ID int `json:"name"`
	Name string `json:"id"`
}

func GetCategories(ctx *iris.Context, storage *story.StoryStorageEngine) ([]Category, error){
	var categories = make([]Category, 2)
	categories[0] = Category{ID:1, Name:"сказки"}
	categories[1] = Category{ID:2, Name:"колыбельные"}
	return categories, nil
}

func getFileForm(ctx *iris.Context, str string) (multipart.File, error){
	info, err := ctx.FormFile(str)
	if(err != nil){
		return  nil, err
	}
	file, err := info.Open()
	if(err != nil){
		return  nil, err
	}
	return file, nil
}

func AddStory(ctx *iris.Context, storage *story.StoryStorageEngine) (int, error) {
	var story_obj = story.Story{}
	err := json.Unmarshal(ctx.Request.Body(), &story_obj)
	log.Print(story_obj)
	if err != nil {
		return 0, err
	}
	id,err := (*storage).Create(story_obj)
	return id, err
}

func AddStoryFile(ctx *iris.Context, id string, googlestorage *storage.Client) bool {
	// Get the file from the request
	audio, err := getFileForm(ctx, "file")
	if err != nil{
		fmt.Println(err)
		return false
	}
	defer audio.Close()


	out1, err := os.OpenFile("./static/audio/"+string(id)+".mp3", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer out1.Close()
	io.Copy(out1, audio)
	dir1 :=  "static/audio/"+id+".raw"
	dir2 :=  "static/audio/"+id+".mp3"
	cmd := exec.Command("mpg123","-O", dir1, "--rate", "8000",  "--mono", "-e", "u8", dir2)
	log.Print(cmd.Args)

	asd, err := cmd.CombinedOutput()
	//asd, err = exec.Command("pwd").CombinedOutput()
	log.Print(string(asd))
	log.Print(err)
	storybckt := (*googlestorage).Bucket("hardteddy_stories")
	if(err != nil){

	}
	file, err := os.Open(dir1)
	defer file.Close()
	obj := storybckt.Object(id+".raw")
	w := obj.NewWriter(context.Background())
	w.ContentType = "audio/basic"
	//w.ACL = []storage.ACLRule{{storage.AllUsers, storage.RoleReader}}
	defer w.Close()
	buf := make([]byte, 2048*16)
	for {
		len, err := file.Read(buf)
		w.Write(buf[0:len])
		if err != nil {
			break
		}
	}
	return true
}

func AddStorySmallImg(ctx *iris.Context, id string, googlestorage *storage.Client) bool {
	// Get the file from the request
	small_img, err := getFileForm(ctx, "file")
	if err != nil{
		fmt.Println(err)
		return false
	}
	defer small_img.Close()
	storybckt := (*googlestorage).Bucket("hardteddy_images")
	if(err != nil){
		log.Print(err)
		return false
	}
	obj := storybckt.Object("small/"+id+".jpg")
	w := obj.NewWriter(context.Background())
	w.ContentType = "image/jpeg"
	w.ACL = []storage.ACLRule{{storage.AllUsers, storage.RoleReader}}
	defer w.Close()
	buf := make([]byte, 2048*16)
	for {
		len, err := small_img.Read(buf)
		w.Write(buf[0:len])
		if err != nil {
			break
		}
	}
	return true
}

func AddStoryLargeImg(ctx *iris.Context, id string, googlestorage *storage.Client) bool {
	large_img, err := getFileForm(ctx, "file")
	if err != nil{
		fmt.Println(err)
		return false
	}
	defer large_img.Close()
	storybckt := (*googlestorage).Bucket("hardteddy_images")
	if(err != nil){
		log.Print(err)
		return false
	}
	obj := storybckt.Object("large/"+id+".jpg")
	w := obj.NewWriter(context.Background())
	w.ContentType = "image/jpeg"
	w.ACL = []storage.ACLRule{{storage.AllUsers, storage.RoleReader}}
	defer w.Close()
	buf := make([]byte, 2048*16)
	for {
		len, err := large_img.Read(buf)
		w.Write(buf[0:len])
		if err != nil {
			break
		}
	}
	return true
}
func AddStoryFiles(ctx *iris.Context, id string) bool {
	// Get the file from the request
	audio, err := getFileForm(ctx, "audio")
	if err != nil{
		fmt.Println(err)
		return false
	}
	defer audio.Close()

	small_img, err := getFileForm(ctx, "small_img")
	if err != nil{
		fmt.Println(err)
		return false
	}
	defer small_img.Close()

	large_img, err := getFileForm(ctx, "large_img")
	if err != nil{
		fmt.Println(err)
		return false
	}
	defer large_img.Close()

	out1, err := os.OpenFile("./static/audio/"+string(id)+".mp3", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer out1.Close()
	out2, err := os.OpenFile("./uploads/large_"+string(id)+".jpg", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer out2.Close()
	out3, err := os.OpenFile("./uploads/small_"+string(id)+".jpg", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer out3.Close()
	io.Copy(out1, audio)
	io.Copy(out2, large_img)
	io.Copy(out3, small_img)


	dir1 :=  "/home/ihelos/Desktop/go/src/github.com/iHelos/tech_teddy/static/audio/test.raw"
	dir2 :=  "/home/ihelos/Desktop/go/src/github.com/iHelos/tech_teddy/static/audio/1.mp3"
	cmd := exec.Command("mpg123","-O",dir1, "--rate", "8000",  "--mono", "-e", "u8", dir2)
	log.Print(cmd.Args)

	asd, err := cmd.CombinedOutput()
	//asd, err = exec.Command("pwd").CombinedOutput()
	log.Print(string(asd))
	log.Print(err)
	return true
}