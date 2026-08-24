package handler

import (
	"net/http"

	"logaggregation/internal/model"
	"logaggregation/pkg/httpx"
)

func (s *Server) registerAlertRuleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/alert-rules", s.createAlertRule)
	mux.HandleFunc("GET /api/alert-rules", s.listAlertRules)
	mux.HandleFunc("GET /api/alert-rules/{id}", s.getAlertRule)
	mux.HandleFunc("PUT /api/alert-rules/{id}", s.updateAlertRule)
	mux.HandleFunc("DELETE /api/alert-rules/{id}", s.deleteAlertRule)
	mux.HandleFunc("POST /api/alert-rules/{id}/disable", s.disableAlertRule)
	mux.HandleFunc("POST /api/alert-rules/{id}/enable", s.enableAlertRule)
	mux.HandleFunc("POST /api/alert-rules/batch-status", s.batchUpdateAlertRuleStatus)
}

type createAlertRuleRequest struct {
	StreamID string `json:"stream_id"`
	Level    string `json:"level"`
	Pattern  string `json:"pattern"`
	Status   string `json:"status"`
}

func (s *Server) createAlertRule(w http.ResponseWriter, r *http.Request) {
	var req createAlertRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.CreateAlertRule(model.AlertRule{StreamID: req.StreamID, Level: req.Level, Pattern: req.Pattern, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, a)
}

func (s *Server) listAlertRules(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.AlertRuleFilter{
		Status:   r.URL.Query().Get("status"),
		StreamID: r.URL.Query().Get("stream_id"),
		Level:    r.URL.Query().Get("level"),
		Keyword:  r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListAlertRules(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getAlertRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := s.svc.GetAlertRule(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

type updateAlertRuleRequest struct {
	StreamID string `json:"stream_id"`
	Level    string `json:"level"`
	Pattern  string `json:"pattern"`
}

func (s *Server) updateAlertRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateAlertRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.UpdateAlertRule(id, model.AlertRule{StreamID: req.StreamID, Level: req.Level, Pattern: req.Pattern})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

func (s *Server) deleteAlertRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteAlertRule(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) disableAlertRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := s.svc.TransitionAlertRuleStatus(id, model.AlertRuleStatusDisabled)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

func (s *Server) enableAlertRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := s.svc.TransitionAlertRuleStatus(id, model.AlertRuleStatusActive)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

type batchUpdateAlertRuleStatusRequest struct {
	IDs    []string `json:"ids"`
	Status string   `json:"status"`
}

func (s *Server) batchUpdateAlertRuleStatus(w http.ResponseWriter, r *http.Request) {
	var req batchUpdateAlertRuleStatusRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if err := s.svc.BatchUpdateAlertRuleStatus(req.IDs, req.Status); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
