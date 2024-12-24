package monitor

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func GetResponse() *Response {
	return &Response{
		Status:  "success",
		Message: "monitor response",
	}
}
