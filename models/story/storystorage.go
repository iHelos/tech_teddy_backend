package story

type StoryStorageEngine interface {
	Create(Story) (int, error)
	Load(string) (Story, error)
	GetAll(order string, order_type string, page int) ([]Story, error)
	GetAllByCategory(order string, order_type string, page int, category int) ([]Story, error)
	GetMyStories(string) ([]Story, error)
	GetSubStories(int)([]SubStory,error)
	Search(keyword string) ([]Story, error)
	SetSizeM(int, int64)
	SetSizeF(int, int64)
	//	UrlMale 	    string	`json:"url_m"`
	//UrlFemale	    string	`json:"url_f"`
	//UrlMp3Male	    string	`json:"url_m_mp3"`
	//UrlMp3Female	    string	`json:"url_f_mp3"`
	//UrlBackground	    string	`json:"url_background"`
	//UrlImageLarge	string		`json:"url_img_large"`
	//UrlImageSmall	string		`json:"url_img_small"`
	SetUrlMale(int, string)
	SetUrlFemale(int, string)
	SetUrlMp3Male(int, string)
	SetUrlMp3Female(int, string)
	SetUrlImageLarge(int, string)
	SetUrlImageSmall(int, string)
	SetUrlBackground(int, string)
}