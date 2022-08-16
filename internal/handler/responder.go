package handler

type ResponseJSON struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func RenderJSON(data interface{}) ResponseJSON {
	var resp ResponseJSON
	switch v := data.(type) {
	case error:
		resp.Message = v.Error()
	default:
		resp.Data = v
	}

	return resp
}
