package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/heppu/miniurl/ui"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/exp/slog"
)

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

type Handler interface {
	AddUrl(url string) (hash string, err error)
	GetUrl(hash string) (url string, err error)
}

type API struct {
	handler Handler
}

func Bind(r *httprouter.Router, h Handler) {
	a := &API{handler: h}
	r.GET("/", a.IndexHandler)
	r.GET("/:hash", a.RedirectHandler)
	r.POST("/api/v1/url", a.PostUrlHandler)
}

func (a *API) RedirectHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	hash := p.ByName("hash")
	url, err := a.handler.GetUrl(hash)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func (a *API) IndexHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	_, err := w.Write(ui.Index)
	if err != nil {
		slog.Error(err.Error())
	}
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

func (a *API) PostUrlHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var req AddUrlReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		resp := ErrorResp{Msg: "bad request"}
		respondJSON(w, http.StatusBadRequest, resp)
		return
	}

	hash, err := a.handler.AddUrl(req.Url)
	if err != nil {
		resp := ErrorResp{Msg: "internal server error"}
		respondJSON(w, http.StatusInternalServerError, resp)
		return
	}

	resp := AddUrlResp{Url: req.Url, Hash: hash}
	respondJSON(w, http.StatusOK, resp)
}

func respondJSON(w http.ResponseWriter, code int, resp any) {
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		slog.Error(err.Error())
	}
}
