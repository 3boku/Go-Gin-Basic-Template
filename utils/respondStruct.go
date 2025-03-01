package utils

type ErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
	Reason string `json:"reason"`
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type GetResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}
