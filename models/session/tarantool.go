package session

import (
	"github.com/tarantool/go-tarantool"
)

type SessionConnection struct {
	*tarantool.Connection
}

func (connection SessionConnection) Load(sid string) map[string]interface{} {
	values := make(map[string]interface{})
	resp, err := connection.Ping()
	if err == nil {
		resp, err = connection.Select("sessions", "primary", 0, 1, tarantool.IterEq, []interface{}{sid})
		if err == nil && len(resp.Data) == 1 {
			data := resp.Data[0].([]interface{})[1] //value withoud sid
			if data != nil {
				for k, v := range data.(map[interface{}]interface{}) {
					values[k.(string)] = v
				}
			}
		}
	}
	return values

}

func (connection SessionConnection) Update(sid string, newValues map[string]interface{}) {
	if len(newValues) == 0 {
		connection.Delete("sessions", "primary", []interface{}{sid})
	} else {
		connection.Replace("sessions", []interface{}{sid, newValues})
	}

}