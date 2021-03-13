package server

import (
	"log"
	"net/http"

	"game-gacha/pkg/http/middleware"
	"game-gacha/pkg/server/handler"
)

func Serve(addr string) {
	http.HandleFunc("/setting/get", get(handler.HandleSettingGet()))
	http.HandleFunc("/user/create", post(handler.HandleUserCreate()))

	http.HandleFunc("/user/get", get(middleware.Authenticate(handler.HandleUserGet())))
	http.HandleFunc("/user/update", post(middleware.Authenticate(handler.HandleUserUpdate())))

	log.Println("Server is running ...")
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Listen and serve failed. %+v", err)
	}
}
func get(api http.HandlerFunc) http.HandlerFunc {
	return httpMethod(api, http.MethodGet)
}
func post(api http.HandlerFunc) http.HandlerFunc {
	return httpMethod(api, http.MethodPost)
}
func httpMethod(api http.HandlerFunc, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// CORS
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Accept, Origin, x-token")
		// preflight request
		if r.Method == http.MethodOptions {
			return
		}
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			if _, err := w.Write([]byte("Method not allowed")); err != nil {
				log.Println(err)
			}
			return
		}
		w.Header().Add("Content-Type", "application/json")
		api(w, r)
	}
}
