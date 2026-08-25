package handler

import (
	"net/http"

	"reverseproxy/internal/model"
	"reverseproxy/pkg/httpx"
)

func (s *Server) registerUpstreamRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/upstreams", s.createUpstream)
	mux.HandleFunc("GET /api/upstreams", s.listUpstream)
	mux.HandleFunc("GET /api/upstreams/{id}", s.getUpstream)
	mux.HandleFunc("PUT /api/upstreams/{id}", s.updateUpstream)
	mux.HandleFunc("DELETE /api/upstreams/{id}", s.deleteUpstream)
	mux.HandleFunc("PATCH /api/upstreams/{id}/status", s.transitionUpstream)
}

type createUpstreamRequest struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Weight   int    `json:"weight"`
	MaxConns int    `json:"max_conns"`
}

func (s *Server) createUpstream(w http.ResponseWriter, r *http.Request) {
	var req createUpstreamRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.CreateUpstream(model.Upstream{Name: req.Name, Host: req.Host, Port: req.Port, Weight: req.Weight, MaxConns: req.MaxConns})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, x)
}

func (s *Server) listUpstream(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.UpstreamFilter{
		Status: r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListUpstreams(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getUpstream(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	x, err := s.svc.GetUpstream(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}

type updateUpstreamRequest struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Weight   int    `json:"weight"`
	MaxConns int    `json:"max_conns"`
}

func (s *Server) updateUpstream(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateUpstreamRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.UpdateUpstream(id, model.Upstream{Name: req.Name, Host: req.Host, Port: req.Port, Weight: req.Weight, MaxConns: req.MaxConns})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}

func (s *Server) deleteUpstream(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteUpstream(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type transitionUpstreamRequest struct {
	Status string `json:"status"`
}

func (s *Server) transitionUpstream(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req transitionUpstreamRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.TransitionUpstream(id, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}
