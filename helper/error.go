package helper

import "encoding/json"

type TeddyError struct {
	Messages map[string][]int
}

func (error TeddyError) Error() string{
	errmsg, _ := json.Marshal(error.Messages)
	return string(errmsg)
}

func NewError() *TeddyError{
	usererror := TeddyError{
		Messages: make(map[string][]int),
	}
	return &usererror
}

func (error TeddyError) Append(key string, value int){
	error.Messages[key] = append(error.Messages[key], value)
}
