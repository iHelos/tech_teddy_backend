package story

//import "encoding/json"

type Story struct  {
	Name string `json:"name"`
	Description string `json:"description"`
	Author string `json:"author"`
	ID string `json:"id"`
	Price int `json:"price"`
}

func GetAllStories() []Story{
	story1 := Story{Name:"Story1", Description:"Story1 awesome description", Author:"iHelos", ID:"1", Price:15}
	story2 := Story{Name:"Story2", Description:"Story2 awesome description", Author:"AnnJelly", ID:"2", Price:25}
	//json1, _ := json.Marshal(story1)
	//json2, _ := json.Marshal(story2)
	return []Story{story1, story2}
}

func GetMyStories() []Story{
	story1 := Story{Name:"Story1", Description:"Story1 awesome description", Author:"iHelos", ID:"1", Price:15}
	//json1, _ := json.Marshal(story1)
	//json2, _ := json.Marshal(story2)
	return []Story{story1}
}