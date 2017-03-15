package helper

func GetResponse(status int, body interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	result["status"] = status
	result["body"] = body
	return result
}