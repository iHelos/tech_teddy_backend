package user

import "encoding/json"

type UserError struct {
	Messages map[string][]int
}

func (error UserError) Error() string{
	errmsg, _ := json.Marshal(error.Messages)
	return string(errmsg)
}

func NewUserError() *UserError{
	usererror := UserError{
		Messages: make(map[string][]int),
	}
	return &usererror
}

func (error UserError) Append(key string, value int){
	error.Messages[key] = append(error.Messages[key], value)
}