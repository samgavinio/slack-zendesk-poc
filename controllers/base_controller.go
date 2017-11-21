package controllers

type (
	Controller struct {}

	JsonResponse struct {
		Status string `json:"status"`
		Message string `json:"message"`
		Data interface{} `json:"data"`
	}
)

func (c *Controller) Success(data interface{}, message string) JsonResponse {
	var empty struct{}
	if data == nil {
		data = empty
	}

	return JsonResponse{Status: "OK", Data: data, Message: message}
}

func (c *Controller) Error(data interface{}, message string) JsonResponse {
	var empty struct{}
	if data == nil {
		data = empty
	}

	return JsonResponse{Status: "Error", Data: data, Message: message}
}
