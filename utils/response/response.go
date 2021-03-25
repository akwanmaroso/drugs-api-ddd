package response

func ResponseJSON(isError bool, status int, msg interface{}) map[string]interface{} {
	return map[string]interface{}{
		"error": isError,
		"status": status,
		"message": msg,
	}
}
