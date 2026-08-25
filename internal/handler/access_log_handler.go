package handler

import (
	"net/http"

	"reverseproxy/internal/model"
	"reverseproxy/pkg/httpx"
)

func (s *Server) registerAccessLogRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/access-logs", s.createAccessLog)
	mux.HandleFunc("GET /api/access-logs", s.listAccessLog)
	mux.HandleFunc("GET /api/access-logs/{id}", s.getAccessLog)
	mux.HandleFunc("DELETE /api/access-logs/{id}", s.deleteAccessLog)
}

type createAccessLogRequest struct {
	RouteID    string `json:"route_id"`
	UpstreamID string `json:"upstream_id"`
	ClientIP   string `json:"client_ip"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	StatusCode int    `json:"status_code"`
	LatencyMs  int64  `json:"latency_ms"`
	Timestamp  int64  `json:"timestamp"`
}

func (s *Server) createAccessLog(w http.ResponseWriter, r *http.Request) {
	var req createAccessLogRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.CreateAccessLog(model.AccessLog{RouteID: req.RouteID, UpstreamID: req.UpstreamID, ClientIP: req.ClientIP, Method: req.Method, Path: req.Path, StatusCode: req.StatusCode, LatencyMs: req.LatencyMs, Timestamp: req.Timestamp})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, x)
}

func (s *Server) listAccessLog(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.AccessLogFilter{
		RouteId:    r.URL.Query().Get("route_id"),
		StatusCode: parseIntQuery(r, "status_code"),
		Method:     r.URL.Query().Get("method"),
	}
	items, total, err := s.svc.ListAccessLogs(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getAccessLog(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	x, err := s.svc.GetAccessLog(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}

func (s *Server) deleteAccessLog(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteAccessLog(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
