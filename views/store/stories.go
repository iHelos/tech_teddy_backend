package store

import (
	"github.com/kataras/iris"
	"github.com/iHelos/tech_teddy/models/story"
	"github.com/iHelos/tech_teddy/models/user"
	"github.com/labstack/gommon/log"
	"strconv"
	"strings"
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