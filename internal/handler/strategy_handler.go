package handler

import (
	"net/http"

	"reverseproxy/internal/model"
	"reverseproxy/pkg/httpx"
)

func (s *Server) registerStrategyRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/strategies", s.createStrategy)
	mux.HandleFunc("GET /api/strategies", s.listStrategy)
	mux.HandleFunc("GET /api/strategies/{id}", s.getStrategy)
	mux.HandleFunc("PUT /api/strategies/{id}", s.updateStrategy)
	mux.HandleFunc("DELETE /api/strategies/{id}", s.deleteStrategy)
	mux.HandleFunc("PATCH /api/strategies/{id}/status", s.transitionStrategy)
}

type createStrategyRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

func (s *Server) createStrategy(w http.ResponseWriter, r *http.Request) {
	var req createStrategyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.CreateStrategy(model.Strategy{Name: req.Name, Type: req.Type, Description: req.Description})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, x)
}

func (s *Server) listStrategy(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.StrategyFilter{
		Type:   r.URL.Query().Get("type"),
		Status: r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListStrategys(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getStrategy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	x, err := s.svc.GetStrategy(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}

type updateStrategyRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

func (s *Server) updateStrategy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateStrategyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.UpdateStrategy(id, model.Strategy{Name: req.Name, Type: req.Type, Description: req.Description})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}

func (s *Server) deleteStrategy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteStrategy(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type transitionStrategyRequest struct {
	Status string `json:"status"`
}

func (s *Server) transitionStrategy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req transitionStrategyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.TransitionStrategy(id, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}
