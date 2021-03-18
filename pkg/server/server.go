package server

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"game-gacha/pkg/db"
	"game-gacha/pkg/http/middleware"
	"game-gacha/pkg/http/response"
	"game-gacha/pkg/server/handler"
	"game-gacha/pkg/server/model"
	"game-gacha/pkg/server/service"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Serve(addr string) {
	httpResponse := response.NewHttpResponse()
	userRepository := model.NewUserRepository(db.Conn)
	mid := middleware.NewMiddleware(httpResponse, userRepository)

	gachaProbabilityRepository := model.NewGachaProbabilityRepository(db.Conn)
	userCollectionItemRepository := model.NewUserCollectionItemRepository(db.Conn)
	collectionItemRepository := model.NewCollectionItemRepository(db.Conn)

	userService := service.NewUserService(userRepository)
	gameService := service.NewGameService(userRepository)
	gachaService := service.NewGachaService(
		userRepository,
		gachaProbabilityRepository,
		userCollectionItemRepository,
		collectionItemRepository,
	)
	rankingService := service.NewRankingService(userRepository)
	collectionService := service.NewCollectionService(
		userRepository,
		userCollectionItemRepository,
		collectionItemRepository,
	)

	userHandler := handler.NewUserHandler(httpResponse, userService)
	settingHandler := handler.NewSettingHandler(httpResponse)
	gameHandler := handler.NewGameHandler(httpResponse, gameService)
	gachaHandler := handler.NewGachaHandler(httpResponse, gachaService)
	rankingHandler := handler.NewRankingHandler(httpResponse, rankingService)
	collectionHandler := handler.NewCollectionHandler(httpResponse, collectionService)

	http.HandleFunc("/setting/get", get(settingHandler.HandleSettingGet))
	http.HandleFunc("/user/create", post(userHandler.HandleUserCreate))

	http.HandleFunc("/user/get", get(mid.Authenticate(userHandler.HandleUserGet)))
	http.HandleFunc("/user/update", post(mid.Authenticate(userHandler.HandleUserUpdate)))

	http.HandleFunc("/gacha/draw", post(mid.Authenticate(gachaHandler.HandleGachaDraw)))

	http.HandleFunc("/game/finish", post(mid.Authenticate(gameHandler.HandleGameFinish)))

	http.HandleFunc("/collection/list", get(mid.Authenticate(collectionHandler.HandleCollectionList)))

	http.HandleFunc("/ranking/list", get(mid.Authenticate(rankingHandler.HandleRankingList)))

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
