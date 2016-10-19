package store

import (
	"github.com/kataras/iris"
	"github.com/iHelos/tech_teddy/models/story"
	"github.com/iHelos/tech_teddy/models/user"
	"github.com/labstack/gommon/log"
)

func GetAllStories(ctx *iris.Context, storage *story.StoryStorageEngine) ([]story.Story, error) {
	story1 := story.Story{Name:"Story1", Description:"Story1 awesome description", Author:"iHelos", ID:1, Price:15}
	story2 := story.Story{Name:"Story2", Description:"Story2 awesome description", Author:"AnnJelly", ID:2, Price:25}
	//json1, _ := json.Marshal(story1)
	//json2, _ := json.Marshal(story2)
	return []story.Story{story1, story2}, nil
}

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