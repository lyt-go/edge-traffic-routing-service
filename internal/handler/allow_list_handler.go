package handler

import (
	"net/http"

	"reverseproxy/internal/model"
	"reverseproxy/pkg/httpx"
)

func (s *Server) registerAllowListRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/allow-lists", s.createAllowList)
	mux.HandleFunc("GET /api/allow-lists", s.listAllowList)
	mux.HandleFunc("GET /api/allow-lists/{id}", s.getAllowList)
	mux.HandleFunc("PUT /api/allow-lists/{id}", s.updateAllowList)
	mux.HandleFunc("DELETE /api/allow-lists/{id}", s.deleteAllowList)
	mux.HandleFunc("PATCH /api/allow-lists/{id}/status", s.transitionAllowList)
}

type createAllowListRequest struct {
	CIDR        string `json:"cidr"`
	Rule        string `json:"rule"`
	Description string `json:"description"`
}

func (s *Server) createAllowList(w http.ResponseWriter, r *http.Request) {
	var req createAllowListRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.CreateAllowList(model.AllowList{CIDR: req.CIDR, Rule: req.Rule, Description: req.Description})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, x)
}

func (s *Server) listAllowList(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.AllowListFilter{
		Rule:   r.URL.Query().Get("rule"),
		Status: r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListAllowLists(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getAllowList(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	x, err := s.svc.GetAllowList(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}

type updateAllowListRequest struct {
	CIDR        string `json:"cidr"`
	Rule        string `json:"rule"`
	Description string `json:"description"`
}

func (s *Server) updateAllowList(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateAllowListRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.UpdateAllowList(id, model.AllowList{CIDR: req.CIDR, Rule: req.Rule, Description: req.Description})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}

func (s *Server) deleteAllowList(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteAllowList(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type transitionAllowListRequest struct {
	Status string `json:"status"`
}

func (s *Server) transitionAllowList(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req transitionAllowListRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.TransitionAllowList(id, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}
