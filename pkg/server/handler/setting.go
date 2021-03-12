package handler

import (
	"net/http"

	"game-gacha/pkg/constant"
	"game-gacha/pkg/http/response"
)

type SettingGetResponse struct{}
type settingGetResponse struct {
	GachaCoinConsumption int `json:"gachaCoinConsumption"`
}

func HandleSettingGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response.Success(w, &settingGetResponse{
			GachaCoinConsumption: constant.GachaCoinConsumption,
		})
	}
}
