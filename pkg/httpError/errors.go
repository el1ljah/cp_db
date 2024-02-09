package httpError

import 	(
	
	"net/http"
)


// NewError example
func NewError(w http.ResponseWriter, error string, code int) {
	er := HTTPError{
		Code:    code,
		Message: error,
	}
	http.Error(w, er.Message, er.Code)
}

// HTTPError example
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}