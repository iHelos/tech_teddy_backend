package tarantool_story_storage

import (
	"github.com/tarantool/go-tarantool"
	"github.com/iHelos/tech_teddy/models/story"
	"encoding/json"
	"log"
)

type StorageConnection struct {
	*tarantool.Connection
}

func (con StorageConnection) Create(st story.Story) (int, error) {
	story_obj, err := con.Call("addStory",
		[]interface{}{
			st.Name,
			st.Description,
			st.Author,
			"00:00",
			st.Price,
			st.SizeM,
			st.SizeF,
			st.UrlMale,
			st.UrlFemale,
			st.UrlMp3Male,
			st.UrlMp3Female,
			st.UrlBackground,
			st.UrlImageLarge,
			st.UrlImageSmall,
		})
	if err!=nil{
		return 0, err
	}
	story_des, err := DeserializeStory(story_obj.Data[0])
	log.Print(story_des, err)
	return int(story_des.ID), err
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

func (con StorageConnection) GetSubStories(id int)(substories []story.SubStory, err error){
	answer, err := con.Call("getSubStories", []interface{}{id})
	if err!= nil{
		return substories,err
	}
	json.Unmarshal([]byte(answer.Data[0].([]interface{})[0].(string)), &substories)
	return substories, err
}
func (con StorageConnection) Search(keyword string) ([]story.Story, error) {
	answer, err := con.Call("findStory", []interface{}{keyword})
	stories, err := DeserializeStoryArray(answer)
	return stories, err
}

func (con StorageConnection) SetSizeM(id int, sizeM int64){
	_,err := con.Update("audio", "primary", []interface{}{uint(id)}, []interface{}{[]interface{}{"=", 7, sizeM}})
	if (err!=nil){
		log.Print(err)
	}
}
func (con StorageConnection) SetSizeF(id int, sizeF int64){
	_,err := con.Update("audio", "primary", []interface{}{uint(id)}, []interface{}{[]interface{}{"=", 8, sizeF}})
	if (err!=nil){
		log.Print(err)
	}
}
func (con StorageConnection) SetUrlMale(id int, url string){
	_,err := con.Update("audio", "primary", []interface{}{uint(id)}, []interface{}{[]interface{}{"=", 9, url}})
	if (err!=nil){
		log.Print(err)
	}
}
func (con StorageConnection) SetUrlFemale(id int, url string){
	_,err := con.Update("audio", "primary", []interface{}{uint(id)}, []interface{}{[]interface{}{"=", 10, url}})
	if (err!=nil){
		log.Print(err)
	}
}
func (con StorageConnection) SetUrlMp3Male(id int, url string){
	_,err := con.Update("audio", "primary", []interface{}{uint(id)}, []interface{}{[]interface{}{"=", 11, url}})
	if (err!=nil){
		log.Print(err)
	}
}
func (con StorageConnection) SetUrlMp3Female(id int, url string){
	_,err := con.Update("audio", "primary", []interface{}{uint(id)}, []interface{}{[]interface{}{"=", 12, url}})
	if (err!=nil){
		log.Print(err)
	}
}
func (con StorageConnection) SetUrlImageLarge(id int, url string){
	_,err := con.Update("audio", "primary", []interface{}{uint(id)}, []interface{}{[]interface{}{"=", 14, url}})
	if (err!=nil){
		log.Print(err)
	}
}
func (con StorageConnection) SetUrlImageSmall(id int, url string){
	_,err := con.Update("audio", "primary", []interface{}{uint(id)}, []interface{}{[]interface{}{"=", 15, url}})
	if (err!=nil){
		log.Print(err)
	}
}
func (con StorageConnection) SetUrlBackground(id int, url string){
	_,err := con.Update("audio", "primary", []interface{}{uint(id)}, []interface{}{[]interface{}{"=", 13, url}})
	if (err!=nil){
		log.Print(err)
	}
}