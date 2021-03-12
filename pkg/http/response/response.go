package response

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Success(w http.ResponseWriter, r interface{}) {
	if r == nil {
		return
	}
	data, err := json.Marshal(r)
	if err != nil {
		log.Println(err)
		InternalServerError(w, "marshal error")
		return
	}
	if _, err = w.Write(data); err != nil {
		log.Println(err)
	}
}

func BadRequest(w http.ResponseWriter, message string) {
	httpError(w, http.StatusBadRequest, message)
}

func InternalServerError(w http.ResponseWriter, message string) {
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
