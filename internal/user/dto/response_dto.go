package dto

type ErrorUnauthorized struct {
	Error string `json:"error" default:"Unauthorized"`
}

type ErrorForbidden struct {
	Error string `json:"error" default:"Forbidden"`
}

type ErrorBadRequest struct {
	Error string `json:"error" default:"Bad Request"`
}

type ErrorNotFound struct {
	Error string `json:"error" default:"Not Found"`
}

type ErrorInternalServer struct {
	Error string `json:"error" default:"Internal Server Error"`
}
