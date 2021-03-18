package handler

import (
	"net/http"

	"game-gacha/pkg/constant"
	"game-gacha/pkg/http/response"
)

type settingGetRequest struct{}
type settingGetResponse struct {
	GachaCoinConsumption int `json:"gachaCoinConsumption"`
}
type settingHandler struct {
	HttpResponse response.HttpResponseInterface
}
type SettingHandlerInterface interface {
	HandleSettingGet(w http.ResponseWriter, r *http.Request)
}

var _ SettingHandlerInterface = (*settingHandler)(nil)

func NewSettingHandler(httpResponse response.HttpResponseInterface) SettingHandlerInterface {
	return &settingHandler{
		HttpResponse: httpResponse,
	}
}

func (h *settingHandler) HandleSettingGet(w http.ResponseWriter, r *http.Request) {
	h.HttpResponse.Success(w, &settingGetResponse{
		GachaCoinConsumption: constant.GachaCoinConsumption,
	})
}
