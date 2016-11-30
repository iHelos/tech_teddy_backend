package tarantool_story_storage

import (
	"github.com/tarantool/go-tarantool"
	"github.com/iHelos/tech_teddy/models/story"
	"errors"
	"strings"
	"strconv"
	"github.com/labstack/gommon/log"
)

func DeserializeStoryArray(response *tarantool.Response) ([]story.Story, error) {
	datalen := len(response.Data)
	if datalen == 1{
		if len(response.Data[0].([]interface{})) == 0{
			return []story.Story{}, nil
		}
	}
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
		//if storyarr, ok := temp[0].([]interface{}); ok {
			if len(temp) >= 16 {
				storyobj, err := arrayToStory(temp)
				return storyobj, err
			}
		//}
	}
	return story.Story{}, errors.New("bad deserialize")
}

func arrayToStory(storyarr []interface{}) (story.Story, error){
	if len(storyarr) >= 16 {
		par1, _ := storyarr[0].(uint64);
		par2, _ := storyarr[1].(uint64);
		par3, _ := storyarr[2].(string);
		par4, _ := storyarr[3].(uint64);
		par5, _ := storyarr[4].(string);
		par6, _ := storyarr[5].(string);
		par7, _ := storyarr[6].(string);
		par8, _ := storyarr[7].(uint64);
		par9, _ := storyarr[8].(uint64);
		par10, _ := storyarr[9].(string);
		par11, _ := storyarr[10].(string);
		par12, _ := storyarr[11].(string);
		par13, _ := storyarr[12].(string);
		par14, _ := storyarr[13].(string);
		par15, _ := storyarr[14].(string);
		par16, _ := storyarr[15].(string);

		//if !(ok1 && ok2 && ok3 && ok4 && ok5 && ok6 && ok7 && ok8 && ok9 && ok10 && ok11 && ok12 && ok13 && ok14 && ok15 && ok16){
		//	log.Print("asdasdasd")
		//	return story.Story{}, errors.New("bad deserialize")
		//}

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
		log.Print("zxc")
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
			SizeM:par8,
			SizeF:par9,
			UrlMale:par10,
			UrlFemale:par11,
			UrlMp3Male:par12,
			UrlMp3Female:par13,
			UrlBackground:par14,
			UrlImageLarge:par15,
			UrlImageSmall:par16,
		}
		return storyobj, nil
	}
	return story.Story{}, errors.New("bad deserialize")
}

