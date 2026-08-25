package handler

import (
	"net/http"

	"reverseproxy/internal/model"
	"reverseproxy/pkg/httpx"
)

func (s *Server) registerHealthCheckRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/health-checks", s.createHealthCheck)
	mux.HandleFunc("GET /api/health-checks", s.listHealthCheck)
	mux.HandleFunc("GET /api/health-checks/{id}", s.getHealthCheck)
	mux.HandleFunc("PUT /api/health-checks/{id}", s.updateHealthCheck)
	mux.HandleFunc("DELETE /api/health-checks/{id}", s.deleteHealthCheck)
	mux.HandleFunc("PATCH /api/health-checks/{id}/status", s.transitionHealthCheck)
}

type createHealthCheckRequest struct {
	UpstreamID    string `json:"upstream_id"`
	Endpoint      string `json:"endpoint"`
	IntervalSec   int    `json:"interval_sec"`
	TimeoutSec    int    `json:"timeout_sec"`
	FailThreshold int    `json:"fail_threshold"`
}

func (s *Server) createHealthCheck(w http.ResponseWriter, r *http.Request) {
	var req createHealthCheckRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.CreateHealthCheck(model.HealthCheck{UpstreamID: req.UpstreamID, Endpoint: req.Endpoint, IntervalSec: req.IntervalSec, TimeoutSec: req.TimeoutSec, FailThreshold: req.FailThreshold})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, x)
}

func (s *Server) listHealthCheck(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.HealthCheckFilter{
		Status:     r.URL.Query().Get("status"),
		UpstreamId: r.URL.Query().Get("upstream_id"),
	}
	items, total, err := s.svc.ListHealthChecks(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getHealthCheck(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	x, err := s.svc.GetHealthCheck(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}

type updateHealthCheckRequest struct {
	UpstreamID    string `json:"upstream_id"`
	Endpoint      string `json:"endpoint"`
	IntervalSec   int    `json:"interval_sec"`
	TimeoutSec    int    `json:"timeout_sec"`
	FailThreshold int    `json:"fail_threshold"`
}

func (s *Server) updateHealthCheck(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateHealthCheckRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.UpdateHealthCheck(id, model.HealthCheck{UpstreamID: req.UpstreamID, Endpoint: req.Endpoint, IntervalSec: req.IntervalSec, TimeoutSec: req.TimeoutSec, FailThreshold: req.FailThreshold})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}

func (s *Server) deleteHealthCheck(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteHealthCheck(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type transitionHealthCheckRequest struct {
	Status string `json:"status"`
}

func (s *Server) transitionHealthCheck(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req transitionHealthCheckRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.TransitionHealthCheck(id, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}
