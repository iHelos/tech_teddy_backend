package model

import (
	"gopkg.in/vmihailenco/msgpack.v2"
	"reflect"
	"fmt"
	"github.com/labstack/gommon/log"
)
type Story struct {
	ID               int
	Category         int
	Name             string
	Price            int
	Duration         string
	Description      string
	AuthorID         int
	Roled            bool
	DurationSplitted Duration
	ImgUrls          UrlImage
	Parts            []StoryPart
}

type StoryPart struct {
	Text string
	Part string
	Audio UrlAudio
}

type Duration struct {
	Minutes int
	Seconds int
}

type UrlImage struct {
	Small string
	Large string
}

type UrlAudio struct {
	Raw 	 string
	Original string
}
//1) ID 		int
//1.1) Category
//2) Name 		string
//2.1) Price
//2.2) DurationStr
//3) Description 	string
//4) AuthorID  		int
//5) Roled		bool
//6) DurationSplitted 		DurationSplitted
//7) ImgUrls 		UrlImage
//8) Parts 		[]StoryPar
//category, name, price, duration, descriprion, author
func encodeStory(e *msgpack.Encoder, v reflect.Value) error {
	m := v.Interface().(Story)
	if m.ID > 0 {
		if err := e.EncodeArrayLen(11); err != nil {
			return err
		}
		if err := e.EncodeInt(m.ID); err != nil {
			return err
		}
	} else{
		if err := e.EncodeArrayLen(10); err != nil {
			return err
		}
	}
	if err := e.EncodeInt(m.Category); err != nil {
		return err
	}
	if err := e.EncodeString(m.Name); err != nil {
		return err
	}
	if err := e.EncodeInt(m.Price); err != nil {
		return err
	}
	if err := e.EncodeString(m.Duration); err != nil {
		return err
	}
	if err := e.EncodeString(m.Description); err != nil {
		return err
	}
	if err := e.EncodeInt(m.AuthorID); err != nil {
		return err
	}
	if err := e.EncodeBool(m.Roled); err != nil {
		return err
	}
	e.Encode(m.DurationSplitted)
	e.Encode(m.ImgUrls)
	if err := e.EncodeArrayLen(len(m.Parts)); err != nil {
		return err
	}
	for _, c := range m.Parts {
		e.Encode(c)
	}
	return nil
}
func decodeStory(d *msgpack.Decoder, v reflect.Value) error {
	var err error
	var l int
	m := v.Addr().Interface().(*Story)
	if l, err = d.DecodeArrayLen(); err != nil {
		return err
	}
	if l != 11 {
		return fmt.Errorf("array len doesn't match: %d", l)
	}
	if m.ID, err = d.DecodeInt(); err != nil {
		return err
	}
	if m.Category, err = d.DecodeInt(); err != nil {
		return err
	}
	if m.Name, err = d.DecodeString(); err != nil {
		return err
	}
	if m.Price, err = d.DecodeInt(); err != nil {
		return err
	}
	if m.Duration, err = d.DecodeString(); err != nil {
		return err
	}
	if m.Description, err = d.DecodeString(); err != nil {
		return err
	}
	if m.AuthorID, err = d.DecodeInt(); err != nil {
		return err
	}
	if m.Roled, err = d.DecodeBool(); err != nil {
		return err
	}
	d.Decode(&m.DurationSplitted)
	d.Decode(&m.ImgUrls)
	if l, err = d.DecodeArrayLen(); err != nil {
		return err
	}
	m.Parts = make([]StoryPart, l)
	for i := 0; i < l; i++ {
		d.Decode(&m.Parts[i])
	}
	return nil
}

//1) Text 	string
//2) Part 	string
//3) Audio 	UrlAudio
func encodeStoryPart(e *msgpack.Encoder, v reflect.Value) error {
	m := v.Interface().(StoryPart)
	if err := e.EncodeArrayLen(3); err != nil {
		return err
	}
	if err := e.EncodeString(m.Text); err != nil {
		return err
	}
	if err := e.EncodeString(m.Part); err != nil {
		return err
	}
	e.Encode(m.Audio)
	return nil
}
func decodeStoryPart(d *msgpack.Decoder, v reflect.Value) error {
	var err error
	var l int
	m := v.Addr().Interface().(*StoryPart)
	if l, err = d.DecodeArrayLen(); err != nil {
		return err
	}
	if l != 3 {
		return fmt.Errorf("array len doesn't match: %d", l)
	}
	if m.Text, err = d.DecodeString(); err != nil {
		return err
	}
	if m.Part, err = d.DecodeString(); err != nil {
		return err
	}
	d.Decode(&m.Audio)
	return nil
}


//1) Minutes int
//2) Seconds int
func encodeDuration(e *msgpack.Encoder, v reflect.Value) error {
	m := v.Interface().(Duration)
	if err := e.EncodeArrayLen(2); err != nil {
		return err
	}
	if err := e.EncodeInt(m.Minutes); err != nil {
		return err
	}
	if err := e.EncodeInt(m.Seconds); err != nil {
		return err
	}
	return nil
}
func decodeDuration(d *msgpack.Decoder, v reflect.Value) error {
	var err error
	var l int
	m := v.Addr().Interface().(*Duration)
	if l, err = d.DecodeArrayLen(); err != nil {
		return err
	}
	if l != 2 {
		return fmt.Errorf("array len doesn't match: %d", l)
	}
	if m.Minutes, err = d.DecodeInt(); err != nil {
		return err
	}
	if m.Seconds, err = d.DecodeInt(); err != nil {
		return err
	}
	return nil
}

//1) Small string
//2) Large string
func encodeUrlImage(e *msgpack.Encoder, v reflect.Value) error {
	m := v.Interface().(UrlImage)
	if err := e.EncodeArrayLen(2); err != nil {
		return err
	}
	if err := e.EncodeString(m.Small); err != nil {
		return err
	}
	if err := e.EncodeString(m.Large); err != nil {
		return err
	}
	return nil
}
func decodeUrlImage(d *msgpack.Decoder, v reflect.Value) error {
	var err error
	var l int
	m := v.Addr().Interface().(*UrlImage)
	if l, err = d.DecodeArrayLen(); err != nil {
		return err
	}
	if l != 2 {
		return fmt.Errorf("array len doesn't match: %d", l)
	}
	if m.Small, err = d.DecodeString(); err != nil {
		return err
	}
	if m.Large, err = d.DecodeString(); err != nil {
		return err
	}
	return nil
}

//1) Raw 	string
//2) Original  	string
func encodeUrlAudio(e *msgpack.Encoder, v reflect.Value) error {
	m := v.Interface().(UrlAudio)
	if err := e.EncodeArrayLen(2); err != nil {
		return err
	}
	if err := e.EncodeString(m.Raw); err != nil {
		return err
	}
	if err := e.EncodeString(m.Original); err != nil {
		return err
	}
	return nil
}
func decodeUrlAudio(d *msgpack.Decoder, v reflect.Value) error {
	var err error
	var l int
	m := v.Addr().Interface().(*UrlAudio)
	if l, err = d.DecodeArrayLen(); err != nil {
		return err
	}
	if l != 2 {
		return fmt.Errorf("array len doesn't match: %d", l)
	}
	if m.Raw, err = d.DecodeString(); err != nil {
		return err
	}
	if m.Original, err = d.DecodeString(); err != nil {
		return err
	}
	return nil
}

func CreateStory(new_story Story) (created_story Story, err error){
	var stories []Story
	err = client.Call17Typed("box.space.audio:auto_increment", []interface{}{new_story}, &stories)
	if (err!=nil){
		log.Print(err)
		return Story{}, err
	}
	return stories[0], err
}
func UpdateStory(new_story Story) (updated_story Story, err error){
	var stories []Story
	//err = client.ReplaceTyped("user", new_profile, &profiles)
	return stories[0], nil
}
func GetStory(id int) (story Story, err error){
	//var profiles []Profile
	//err = client.SelectTyped("user", "primary", 0,1, tarantool.IterEq, []interface{}{uint(id)}, &profiles)
	//if len(profiles)>0 {
	//	return profiles[0], err
	//}
	return Story{}, nil
}
func GetStoriesByUser(user_id int) (story []Story, err error){
	var stories []Story
	err = client.Call17Typed("getUserStories", []interface{}{user_id}, &stories)
	return stories, nil
}

func Search(str string)(story []Story, err error){
	var stories []Story
	err = client.Call17Typed("findStory", []interface{}{str}, &stories)
	return stories, nil
}

const limit int = 15
var order_types map[string]int = map[string]int{
	"desc":0,
	"asc":1,
}
var orders map[string]int = map[string]int{
	"name":0,
	"price":1,
	"duration":2,
}

func GetAll(order string, order_type string, page int) ([]Story, error) {
	offset := limit * page
	var stories [][]Story
	var order_code int = orders[order]
	var order_type_code int = order_types[order_type]
	err := client.Call17Typed("getAllStories", []interface{}{offset, limit, order_code, order_type_code}, &stories)
	if err!=nil {return []Story{}, err}
	return stories[0], err
}
func GetAllByCategory(order string, order_type string, page int, category int) ([]Story, error) {
	offset := limit * page
	var stories [][]Story
	var order_code int = orders[order]
	var order_type_code int = order_types[order_type]
	err := client.Call17Typed("getAllCategoryStories", []interface{}{category, offset, limit, order_code, order_type_code}, &stories)
	if err!=nil {return []Story{}, err}
	return stories[0], err
}