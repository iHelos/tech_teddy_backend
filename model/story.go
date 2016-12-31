package model

import (
	"gopkg.in/vmihailenco/msgpack.v2"
	"reflect"
	"fmt"
	"github.com/tarantool/go-tarantool"
)

type Story struct {
	ID 		int
	Name 		string
	Description 	string
	AuthorID  	int
	isRoled		bool
	Duration 	Duration
	ImgUrls 	UrlImage
	Parts 		[]StoryPart
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

func encodeStory(e *msgpack.Encoder, v reflect.Value) error {
	m := v.Interface().(Story)
	if m.ID > 0 {
		if err := e.EncodeArrayLen(8); err != nil {
			return err
		}
		if err := e.EncodeInt(m.ID); err != nil {
			return err
		}
	} else{
		if err := e.EncodeArrayLen(7); err != nil {
			return err
		}
	}
	if err := e.EncodeString(m.Name); err != nil {
		return err
	}
	if err := e.EncodeString(m.Description); err != nil {
		return err
	}
	if err := e.EncodeString(m.AuthorID); err != nil {
		return err
	}
	if err := e.EncodeArrayLen(len(m.Parts)); err != nil {
		return err
	}
	for _, l := range m.Likes {
		e.EncodeInt(l)
	}
	if err := e.EncodeArrayLen(len(m.Bears)); err != nil {
		return err
	}
	for _, b := range m.Likes {
		e.EncodeInt(b)
	}
	return nil
}

func decodeStory(d *msgpack.Decoder, v reflect.Value) error {
	var err error
	var l int
	m := v.Addr().Interface().(*Profile)
	if l, err = d.DecodeArrayLen(); err != nil {
		return err
	}
	if l != 6 {
		return fmt.Errorf("array len doesn't match: %d", l)
	}
	if m.ID, err = d.DecodeInt(); err != nil {
		return err
	}
	if m.Email, err = d.DecodeString(); err != nil {
		return err
	}
	if m.Name, err = d.DecodeString(); err != nil {
		return err
	}
	if m.HashPassword, err = d.DecodeString(); err != nil {
		return err
	}
	if l, err = d.DecodeArrayLen(); err != nil {
		return err
	}
	m.Likes = make([]int, l)
	for i := 0; i < l; i++ {
		m.Likes[i],_ = d.DecodeInt()
	}
	if l, err = d.DecodeArrayLen(); err != nil {
		return err
	}
	m.Bears = make([]int, l)
	for i := 0; i < l; i++ {
		m.Bears[i],_ = d.DecodeInt()
	}
	return nil
}

func encodeStoryPart(e *msgpack.Encoder, v reflect.Value) error {
	m := v.Interface().(Profile)
	if m.ID > 0 {
		if err := e.EncodeArrayLen(6); err != nil {
			return err
		}
		if err := e.EncodeInt(m.ID); err != nil {
			return err
		}
	} else{
		if err := e.EncodeArrayLen(5); err != nil {
			return err
		}
	}
	if err := e.EncodeString(m.Email); err != nil {
		return err
	}
	if err := e.EncodeString(m.Name); err != nil {
		return err
	}
	if err := e.EncodeString(m.HashPassword); err != nil {
		return err
	}
	if err := e.EncodeArrayLen(len(m.Likes)); err != nil {
		return err
	}
	for _, l := range m.Likes {
		e.EncodeInt(l)
	}
	if err := e.EncodeArrayLen(len(m.Bears)); err != nil {
		return err
	}
	for _, b := range m.Likes {
		e.EncodeInt(b)
	}
	return nil
}

func decodeStoryPart(d *msgpack.Decoder, v reflect.Value) error {
	var err error
	var l int
	m := v.Addr().Interface().(*Profile)
	if l, err = d.DecodeArrayLen(); err != nil {
		return err
	}
	if l != 6 {
		return fmt.Errorf("array len doesn't match: %d", l)
	}
	if m.ID, err = d.DecodeInt(); err != nil {
		return err
	}
	if m.Email, err = d.DecodeString(); err != nil {
		return err
	}
	if m.Name, err = d.DecodeString(); err != nil {
		return err
	}
	if m.HashPassword, err = d.DecodeString(); err != nil {
		return err
	}
	if l, err = d.DecodeArrayLen(); err != nil {
		return err
	}
	m.Likes = make([]int, l)
	for i := 0; i < l; i++ {
		m.Likes[i],_ = d.DecodeInt()
	}
	if l, err = d.DecodeArrayLen(); err != nil {
		return err
	}
	m.Bears = make([]int, l)
	for i := 0; i < l; i++ {
		m.Bears[i],_ = d.DecodeInt()
	}
	return nil
}

func encodeDuration(e *msgpack.Encoder, v reflect.Value) error {
	m := v.Interface().(Profile)
	if m.ID > 0 {
		if err := e.EncodeArrayLen(6); err != nil {
			return err
		}
		if err := e.EncodeInt(m.ID); err != nil {
			return err
		}
	} else{
		if err := e.EncodeArrayLen(5); err != nil {
			return err
		}
	}
	if err := e.EncodeString(m.Email); err != nil {
		return err
	}
	if err := e.EncodeString(m.Name); err != nil {
		return err
	}
	if err := e.EncodeString(m.HashPassword); err != nil {
		return err
	}
	if err := e.EncodeArrayLen(len(m.Likes)); err != nil {
		return err
	}
	for _, l := range m.Likes {
		e.EncodeInt(l)
	}
	if err := e.EncodeArrayLen(len(m.Bears)); err != nil {
		return err
	}
	for _, b := range m.Likes {
		e.EncodeInt(b)
	}
	return nil
}

func decodeDuration(d *msgpack.Decoder, v reflect.Value) error {
	var err error
	var l int
	m := v.Addr().Interface().(*Profile)
	if l, err = d.DecodeArrayLen(); err != nil {
		return err
	}
	if l != 6 {
		return fmt.Errorf("array len doesn't match: %d", l)
	}
	if m.ID, err = d.DecodeInt(); err != nil {
		return err
	}
	if m.Email, err = d.DecodeString(); err != nil {
		return err
	}
	if m.Name, err = d.DecodeString(); err != nil {
		return err
	}
	if m.HashPassword, err = d.DecodeString(); err != nil {
		return err
	}
	if l, err = d.DecodeArrayLen(); err != nil {
		return err
	}
	m.Likes = make([]int, l)
	for i := 0; i < l; i++ {
		m.Likes[i],_ = d.DecodeInt()
	}
	if l, err = d.DecodeArrayLen(); err != nil {
		return err
	}
	m.Bears = make([]int, l)
	for i := 0; i < l; i++ {
		m.Bears[i],_ = d.DecodeInt()
	}
	return nil
}

func encodeUrlImage(e *msgpack.Encoder, v reflect.Value) error {
	m := v.Interface().(Profile)
	if m.ID > 0 {
		if err := e.EncodeArrayLen(6); err != nil {
			return err
		}
		if err := e.EncodeInt(m.ID); err != nil {
			return err
		}
	} else{
		if err := e.EncodeArrayLen(5); err != nil {
			return err
		}
	}
	if err := e.EncodeString(m.Email); err != nil {
		return err
	}
	if err := e.EncodeString(m.Name); err != nil {
		return err
	}
	if err := e.EncodeString(m.HashPassword); err != nil {
		return err
	}
	if err := e.EncodeArrayLen(len(m.Likes)); err != nil {
		return err
	}
	for _, l := range m.Likes {
		e.EncodeInt(l)
	}
	if err := e.EncodeArrayLen(len(m.Bears)); err != nil {
		return err
	}
	for _, b := range m.Likes {
		e.EncodeInt(b)
	}
	return nil
}

func decodeUrlImage(d *msgpack.Decoder, v reflect.Value) error {
	var err error
	var l int
	m := v.Addr().Interface().(*Profile)
	if l, err = d.DecodeArrayLen(); err != nil {
		return err
	}
	if l != 6 {
		return fmt.Errorf("array len doesn't match: %d", l)
	}
	if m.ID, err = d.DecodeInt(); err != nil {
		return err
	}
	if m.Email, err = d.DecodeString(); err != nil {
		return err
	}
	if m.Name, err = d.DecodeString(); err != nil {
		return err
	}
	if m.HashPassword, err = d.DecodeString(); err != nil {
		return err
	}
	if l, err = d.DecodeArrayLen(); err != nil {
		return err
	}
	m.Likes = make([]int, l)
	for i := 0; i < l; i++ {
		m.Likes[i],_ = d.DecodeInt()
	}
	if l, err = d.DecodeArrayLen(); err != nil {
		return err
	}
	m.Bears = make([]int, l)
	for i := 0; i < l; i++ {
		m.Bears[i],_ = d.DecodeInt()
	}
	return nil
}

func encodeUrlAudio(e *msgpack.Encoder, v reflect.Value) error {
	m := v.Interface().(Profile)
	if m.ID > 0 {
		if err := e.EncodeArrayLen(6); err != nil {
			return err
		}
		if err := e.EncodeInt(m.ID); err != nil {
			return err
		}
	} else{
		if err := e.EncodeArrayLen(5); err != nil {
			return err
		}
	}
	if err := e.EncodeString(m.Email); err != nil {
		return err
	}
	if err := e.EncodeString(m.Name); err != nil {
		return err
	}
	if err := e.EncodeString(m.HashPassword); err != nil {
		return err
	}
	if err := e.EncodeArrayLen(len(m.Likes)); err != nil {
		return err
	}
	for _, l := range m.Likes {
		e.EncodeInt(l)
	}
	if err := e.EncodeArrayLen(len(m.Bears)); err != nil {
		return err
	}
	for _, b := range m.Likes {
		e.EncodeInt(b)
	}
	return nil
}

func decodeUrlAudio(d *msgpack.Decoder, v reflect.Value) error {
	var err error
	var l int
	m := v.Addr().Interface().(*Profile)
	if l, err = d.DecodeArrayLen(); err != nil {
		return err
	}
	if l != 6 {
		return fmt.Errorf("array len doesn't match: %d", l)
	}
	if m.ID, err = d.DecodeInt(); err != nil {
		return err
	}
	if m.Email, err = d.DecodeString(); err != nil {
		return err
	}
	if m.Name, err = d.DecodeString(); err != nil {
		return err
	}
	if m.HashPassword, err = d.DecodeString(); err != nil {
		return err
	}
	if l, err = d.DecodeArrayLen(); err != nil {
		return err
	}
	m.Likes = make([]int, l)
	for i := 0; i < l; i++ {
		m.Likes[i],_ = d.DecodeInt()
	}
	if l, err = d.DecodeArrayLen(); err != nil {
		return err
	}
	m.Bears = make([]int, l)
	for i := 0; i < l; i++ {
		m.Bears[i],_ = d.DecodeInt()
	}
	return nil
}

func createStory(new_story Story) (created_story Story, err error){
	var profiles []Profile
	err = client.Call17Typed("box.space.user:auto_increment", []interface{}{new_profile}, &profiles)
	return profiles[0], err
}

func updateStory(new_story Story) (updated_story Story, err error){
	var profiles []Profile
	err = client.ReplaceTyped("user", new_profile, &profiles)
	return profiles[0], err
}

func getStory(id int) (story Story, err error){
	var profiles []Profile
	err = client.SelectTyped("user", "primary", 0,1, tarantool.IterEq, []interface{}{uint(id)}, &profiles)
	if len(profiles)>0 {
		return profiles[0], err
	}
	return profile, errors.New("not exists")
}
