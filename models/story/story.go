package story

//import "encoding/json"

type Story struct {
	Name        string 	`json:"name"`
	Description string 	`json:"description"`
	Author      string 	`json:"author"`
	Category    uint64 	`json:"category"`
	Minutes     int 	`json:"minutes"`
	Seconds     int 	`json:"seconds"`
	ID          uint64 	`json:"id"`
	Price       uint64 	`json:"price"`
	SizeM 	    uint64	`json:"size_m"`
	SizeF	    uint64	`json:"size_f"`
	//Url 	    string	`json:"url"`
	UrlMale 	    string	`json:"url_m"`
	UrlFemale	    string	`json:"url_f"`
	UrlMp3Male	    string	`json:"url_m_mp3"`
	UrlMp3Female	    string	`json:"url_f_mp3"`
	UrlBackground	    string	`json:"url_background"`
	UrlImageLarge	string		`json:"url_img_large"`
	UrlImageSmall	string		`json:"url_img_small"`
}