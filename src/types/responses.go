package types

//GenericResponse - generic http response struct
type GenericResponse struct {
	Response bool `json:"response"`
}

//ErrorResponse - returns an error response
type ErrorResponse struct {
	Error ErrorResponseBody `json:"error"`
}

//ErrorResponseBody - returns an error response
type ErrorResponseBody struct {
	HTTPStatusCode int    `json:"statusCode"`
	ErrorCode      int    `json:"errorCode"`
	ErrorMsg       string `json:"errorMsg"`
}
