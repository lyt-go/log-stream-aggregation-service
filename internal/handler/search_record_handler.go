package handler

import (
	"net/http"

	"logaggregation/internal/model"
	"logaggregation/pkg/httpx"
)

func (s *Server) registerSearchRecordRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/search-records", s.createSearchRecord)
	mux.HandleFunc("GET /api/search-records", s.listSearchRecords)
	mux.HandleFunc("GET /api/search-records/{id}", s.getSearchRecord)
	mux.HandleFunc("PUT /api/search-records/{id}", s.updateSearchRecord)
	mux.HandleFunc("DELETE /api/search-records/{id}", s.deleteSearchRecord)
}

type createSearchRecordRequest struct {
	Query       string `json:"query"`
	ResultCount int    `json:"result_count"`
	DurationMs  int    `json:"duration_ms"`
}

func (s *Server) createSearchRecord(w http.ResponseWriter, r *http.Request) {
	var req createSearchRecordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sr, err := s.svc.CreateSearchRecord(model.SearchRecord{Query: req.Query, ResultCount: req.ResultCount, DurationMs: req.DurationMs})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, sr)
}

func (s *Server) listSearchRecords(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.SearchRecordFilter{
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListSearchRecords(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getSearchRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sr, err := s.svc.GetSearchRecord(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sr)
}

type updateSearchRecordRequest struct {
	Query       string `json:"query"`
	ResultCount int    `json:"result_count"`
	DurationMs  int    `json:"duration_ms"`
}

func (s *Server) updateSearchRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateSearchRecordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sr, err := s.svc.UpdateSearchRecord(id, model.SearchRecord{Query: req.Query, ResultCount: req.ResultCount, DurationMs: req.DurationMs})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sr)
}

func (s *Server) deleteSearchRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteSearchRecord(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
