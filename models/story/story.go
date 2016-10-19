package story

//import "encoding/json"

type Story struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Category    uint64 `json:"category"`
	Minutes     int `json:"minutes"`
	Seconds     int `json:"seconds"`
	ID          uint64 `json:"id"`
	Price       uint64 `json:"price"`
}