package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/exp/slog"
)

type Handler interface {
	AddUrl(original string) (hash string, err error)
}

type API struct {
	handler Handler
}

func Bind(r *httprouter.Router, h Handler) {
	a := &API{handler: h}
	r.POST("/api/v1/url", a.AddUrl)
}

type AddUrlReq struct {
	Url string `json:"url"`
}

type AddUrlResp struct {
	Url  string `json:"url"`
	Hash string `json:"hash"`
}

func (a *API) AddUrl(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var v AddUrlReq
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		// TODO
		return
	}

	hash, err := a.handler.AddUrl(v.Url)
	if err != nil {
		// TODO
		return
	}

	resp := AddUrlResp{Url: v.Url, Hash: hash}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		slog.Error(err.Error())
	}
}
