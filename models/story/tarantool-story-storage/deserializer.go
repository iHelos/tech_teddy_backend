package tarantool_story_storage

import (
	"github.com/tarantool/go-tarantool"
	"github.com/iHelos/tech_teddy/models/story"
	"errors"
	"strings"
	"strconv"
)

func DeserializeStoryArray(response *tarantool.Response) ([]story.Story, error) {
	var slice = make([]story.Story, len(response.Data))
	for i, story := range response.Data {
		var err error
		slice[i], err = DeserializeStory(story)
		if (err!=nil){
			return slice, err
		}

	}
	return slice, nil
}

func DeserializeStory(serstory interface{}) (story.Story, error) {
	if temp, ok := serstory.([]interface{}); ok {
		if storyarr, ok := temp[0].([]interface{}); ok {
			if len(storyarr) == 7 {
				storyobj, err := arrayToStory(storyarr)
				return storyobj, err
			}
		}
	}
	return story.Story{}, errors.New("bad deserialize")
}

func arrayToStory(storyarr []interface{}) (story.Story, error){
	if len(storyarr) == 7 {
		par1, ok1 := storyarr[0].(uint64);
		par2, ok2 := storyarr[1].(uint64);
		par3, ok3 := storyarr[2].(string);
		par4, ok4 := storyarr[3].(uint64);
		par5, ok5 := storyarr[4].(string);
		par6, ok6 := storyarr[5].(string);
		par7, ok7 := storyarr[6].(string);

		if !(ok1 && ok2 && ok3 && ok4 && ok5 && ok6 && ok7){
			return story.Story{}, errors.New("bad deserialize")
		}

		times := strings.Split(par5, ":")
		var minutes int
		var seconds int
		minutes,err := strconv.Atoi(times[0])
		if (err != nil){
			return story.Story{}, err
		}
		seconds,err = strconv.Atoi(times[1])
		if (err != nil){
			return story.Story{}, err
		}
		storyobj := story.Story{
			// --id, category, name, price, duration, descriprion, author
			ID:par1,
			Category:par2,
			Name:par3,
			Price:par4,
			Minutes:minutes,
			Seconds:seconds,
			Description:par6,
			Author:par7,
		}
		return storyobj, nil
	}
	return story.Story{}, errors.New("bad deserialize")
}