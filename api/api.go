package api

import (
	"context"
	_ "embed"
	"encoding/json"
	"net/http"

	"github.com/heppu/miniurl/ui"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/exp/slog"
)

type Handler interface {
	AddUrl(original string) (hash string, err error)
	GetUrl(hash string) (url string, err error)
}

type API struct {
	handler Handler
}

type Server struct {
	srv *http.Server
}

func NewServer(listenAddr string, h Handler) *Server {
	r := httprouter.New()
	Bind(r, h)
	return &Server{
		srv: &http.Server{
			Addr:    listenAddr,
			Handler: r,
		},
	}
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.srv.Shutdown(context.Background())
}

func Bind(r *httprouter.Router, h Handler) {
	a := &API{handler: h}
	r.GET("/", a.Index)
	r.GET("/:hash", a.Redirect)
	r.POST("/api/v1/url", a.AddUrl)
}

type AddUrlReq struct {
	Url string `json:"url"`
}

type AddUrlResp struct {
	Url  string `json:"url"`
	Hash string `json:"hash"`
}

type ErrorResp struct {
	Msg string `json:"msg"`
}

func (a *API) Redirect(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	hash := p.ByName("hash")
	url, err := a.handler.GetUrl(hash)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (a *API) Index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	_, err := w.Write(ui.Index)
	if err != nil {
		slog.Error(err.Error())
	}
}

func (a *API) AddUrl(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var v AddUrlReq
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		resp := ErrorResp{Msg: "bad request"}
		respondJSON(w, http.StatusBadRequest, resp)
		return
	}

	hash, err := a.handler.AddUrl(v.Url)
	if err != nil {
		resp := ErrorResp{Msg: "internal server error"}
		respondJSON(w, http.StatusInternalServerError, resp)
		return
	}

	resp := AddUrlResp{Url: v.Url, Hash: hash}
	respondJSON(w, http.StatusOK, resp)
}

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		slog.Error(err.Error())
	}
}
