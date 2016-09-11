package sessionDB

import (
	"github.com/tarantool/go-tarantool"

)

type SessionConnection struct {
	*tarantool.Connection
}

func (connection SessionConnection) Load(sid string) map[string]interface{} {
	values := make(map[string]interface{})
	//_, err := connection.Ping()
	//if err == nil {
	//	println("Tarantool connection error" + err.Error())
	//}
	resp, err := connection.Select("sessions", "primary", 0, 1, tarantool.IterEq, []interface{}{sid})
	if err == nil && len(resp.Data) == 1 {
		data := resp.Data[0].([]interface{})[1] //value withoud sid
		if data != nil {
			for k, v := range data.(map[interface{}]interface{}) {
				values[k.(string)] = v
			}
		}
	}
	return values

}

// update updates the real redis store
func (connection SessionConnection) Update(sid string, newValues map[string]interface{}) {
	if len(newValues) == 0 {
		go connection.Delete("sessions", "primary", []interface{}{sid})
	} else {
		go connection.Replace("sessions", []interface{}{sid, newValues})
	}

}