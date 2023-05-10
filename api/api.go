package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Handler interface {
	AddUrl(url string) (hash string, err error)
}

type API struct {
	handler Handler
}

func Bind(r *httprouter.Router, h Handler) {
	a := &API{handler: h}
	r.POST("/api/v1/url", a.AddUrl)
}

func (a *API) AddUrl(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}
