package model

type Error struct {
	Error string
}

type ErrOrStatus struct {
	Error Error `json:"message"`
}

// ResponseError ...
type ResponseError struct {
	Error interface{} `json:"error"`
}

//InternalServerError ...
type InternalServerError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
