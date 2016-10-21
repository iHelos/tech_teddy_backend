package tarantool_story_storage

import (
	"github.com/tarantool/go-tarantool"
	"github.com/iHelos/tech_teddy/models/story"
)

type StorageConnection struct {
	*tarantool.Connection
}

func (StorageConnection) Create(story.Story) (error) {
	return nil
}
func (StorageConnection) Load(string) (story.Story, error) {
	var storyobj = story.Story{}
	return storyobj, nil
}

const (
	limit int = 25
)

var order_types map[string]int = map[string]int{
	"desc":0,
	"asc":1,
}
var orders map[string]int = map[string]int{
	"name":0,
	"price":1,
	"duration":2,
}

//getAllStories(offset, limit, order, ordertype)
//getAllCategoryStories(category, offset, limit, order, ordertype)
//order = 0 - имя; 1 - цена; 2 - продолжительность
func (con StorageConnection) GetAll(order string, order_type string, page int) ([]story.Story, error) {
	offset := limit * page
	var order_code int = orders[order]
	var order_type_code int = order_types[order_type]
	answer, err := con.Call("getAllStories", []interface{}{offset, limit, order_code, order_type_code })
	stories, err := DeserializeStoryArray(answer)
	return stories, err
}
func (con StorageConnection) GetAllByCategory(order string, order_type string, page int, category int) ([]story.Story, error) {
	offset := limit * page
	var order_code int = orders[order]
	var order_type_code int = order_types[order_type]
	answer, err := con.Call("getAllCategoryStories", []interface{}{category, offset, limit, order_code, order_type_code })
	stories, err := DeserializeStoryArray(answer)
	return stories, err
}

func (con StorageConnection) GetMyStories(login string) ([]story.Story, error) {
	answer, err := con.Call("getUserStories", []interface{}{login})
	stories, err := DeserializeStoryArray(answer)
	return stories, err
}
func (con StorageConnection) Search(keyword string) ([]story.Story, error) {
	answer, err := con.Call("findStory", []interface{}{keyword})
	stories, err := DeserializeStoryArray(answer)
	return stories, err
}