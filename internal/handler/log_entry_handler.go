package handler

import (
	"net/http"

	"logaggregation/internal/model"
	"logaggregation/pkg/httpx"
)

func (s *Server) registerLogEntryRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/log-entries", s.createLogEntry)
	mux.HandleFunc("GET /api/log-entries", s.listLogEntries)
	mux.HandleFunc("GET /api/log-entries/{id}", s.getLogEntry)
	mux.HandleFunc("PUT /api/log-entries/{id}", s.updateLogEntry)
	mux.HandleFunc("DELETE /api/log-entries/{id}", s.deleteLogEntry)
	mux.HandleFunc("POST /api/log-entries/batch-delete", s.batchDeleteLogEntries)
}

type createLogEntryRequest struct {
	StreamID string   `json:"stream_id"`
	Level    string   `json:"level"`
	Message  string   `json:"message"`
	Tags     []string `json:"tags"`
}

func (s *Server) createLogEntry(w http.ResponseWriter, r *http.Request) {
	var req createLogEntryRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	e, err := s.svc.CreateLogEntry(model.LogEntry{StreamID: req.StreamID, Level: req.Level, Message: req.Message, Tags: req.Tags})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, e)
}

func (s *Server) listLogEntries(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.LogEntryFilter{
		StreamID: r.URL.Query().Get("stream_id"),
		Level:    r.URL.Query().Get("level"),
		Keyword:  r.URL.Query().Get("keyword"),
		Tag:      r.URL.Query().Get("tag"),
	}
	items, total, err := s.svc.ListLogEntries(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getLogEntry(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	e, err := s.svc.GetLogEntry(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}

type updateLogEntryRequest struct {
	StreamID string   `json:"stream_id"`
	Level    string   `json:"level"`
	Message  string   `json:"message"`
	Tags     []string `json:"tags"`
}

func (s *Server) updateLogEntry(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateLogEntryRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	e, err := s.svc.UpdateLogEntry(id, model.LogEntry{StreamID: req.StreamID, Level: req.Level, Message: req.Message, Tags: req.Tags})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}

func (s *Server) deleteLogEntry(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteLogEntry(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type batchDeleteLogEntriesRequest struct {
	IDs []string `json:"ids"`
}

func (s *Server) batchDeleteLogEntries(w http.ResponseWriter, r *http.Request) {
	var req batchDeleteLogEntriesRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if err := s.svc.BatchDeleteLogEntries(req.IDs); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
