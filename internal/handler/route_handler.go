package handler

import (
	"net/http"

	"reverseproxy/internal/model"
	"reverseproxy/pkg/httpx"
)

func (s *Server) registerRouteRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/routes", s.createRoute)
	mux.HandleFunc("GET /api/routes", s.listRoute)
	mux.HandleFunc("GET /api/routes/{id}", s.getRoute)
	mux.HandleFunc("PUT /api/routes/{id}", s.updateRoute)
	mux.HandleFunc("DELETE /api/routes/{id}", s.deleteRoute)
	mux.HandleFunc("PATCH /api/routes/{id}/status", s.transitionRoute)
}

type createRouteRequest struct {
	Path        string   `json:"path"`
	Method      string   `json:"method"`
	UpstreamIDs []string `json:"upstream_ids"`
	Priority    int      `json:"priority"`
	StripPrefix bool     `json:"strip_prefix"`
}

func (s *Server) createRoute(w http.ResponseWriter, r *http.Request) {
	var req createRouteRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.CreateRoute(model.Route{Path: req.Path, Method: req.Method, UpstreamIDs: req.UpstreamIDs, Priority: req.Priority, StripPrefix: req.StripPrefix})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, x)
}

func (s *Server) listRoute(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RouteFilter{
		Method: r.URL.Query().Get("method"),
		Status: r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListRoutes(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRoute(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	x, err := s.svc.GetRoute(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}

type updateRouteRequest struct {
	Path        string   `json:"path"`
	Method      string   `json:"method"`
	UpstreamIDs []string `json:"upstream_ids"`
	Priority    int      `json:"priority"`
	StripPrefix bool     `json:"strip_prefix"`
}

func (s *Server) updateRoute(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateRouteRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.UpdateRoute(id, model.Route{Path: req.Path, Method: req.Method, UpstreamIDs: req.UpstreamIDs, Priority: req.Priority, StripPrefix: req.StripPrefix})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}

func (s *Server) deleteRoute(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteRoute(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type transitionRouteRequest struct {
	Status string `json:"status"`
}

func (s *Server) transitionRoute(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req transitionRouteRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.TransitionRoute(id, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}
