package tarantool_story_storage

import (
	"github.com/tarantool/go-tarantool"
	"github.com/iHelos/tech_teddy/models/story"
)

type StorageConnection struct {
	*tarantool.Connection
}

func (StorageConnection) Create(story.Story) (error){
	return nil
}
func (StorageConnection) Load(string) (story.Story, error){
	var storyobj = story.Story{}
	return storyobj, nil
}
func (StorageConnection) GetAll(category int, order string, page int) ([]story.Story, error){
	return []story.Story{}, nil
}
func (con StorageConnection) GetMyStories(login string) ([]story.Story, error){
	answer, err := con.Call("getUserStories", []interface{}{login})
	stories, err := DeserializeStoryArray(answer)
	return stories, err
}
func (StorageConnection) Search(keyword string) ([]story.Story, error){
	return []story.Story{}, nil
}