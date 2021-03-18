package response

import (
	"encoding/json"
	"log"
	"net/http"

	"game-gacha/pkg/derror"
)

type httpResponse struct{}
type HttpResponseInterface interface {
	Success(w http.ResponseWriter, r interface{})
	Failed(w http.ResponseWriter, message string, err error, code int)
}

func NewHttpResponse() HttpResponseInterface {
	return &httpResponse{}
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (hr *httpResponse) Success(w http.ResponseWriter, r interface{}) {
	if r == nil {
		return
	}
	data, err := json.Marshal(r)
	if err != nil {
		log.Println(err)
		internalServerError(w, "marshal error")
		return
	}
	if _, err = w.Write(data); err != nil {
		log.Println(err)
	}
}
func (hr *httpResponse) Failed(w http.ResponseWriter, message string, err error, code int) {
	e := derror.ApplicationError{
		Message: message,
		Err:     err,
		Code:    code,
	}
	log.Println(e)
	switch e.Code {
	case http.StatusBadRequest:
		badRequest(w, "Bad Request")
	case http.StatusInternalServerError:
		internalServerError(w, "Internal Server Error")
	default:
		internalServerError(w, "Unknown Internal Server Error")
	}
}

func badRequest(w http.ResponseWriter, message string) {
	httpError(w, http.StatusBadRequest, message)
}

func internalServerError(w http.ResponseWriter, message string) {
	httpError(w, http.StatusInternalServerError, message)
}
func httpError(w http.ResponseWriter, code int, message string) {
	data, err := json.Marshal(&errorResponse{
		Code:    code,
		Message: message,
	})
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(code)
	if data != nil {
		if _, err = w.Write(data); err != nil {
			log.Println(err)
		}
	}
}
