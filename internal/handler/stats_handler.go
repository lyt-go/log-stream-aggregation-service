package handler

import (
	"net/http"
	"strconv"
	"time"

	"logaggregation/pkg/httpx"
)

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/by-stream", s.statsByStream)
	mux.HandleFunc("GET /api/stats/by-level", s.statsByLevel)
	mux.HandleFunc("GET /api/stats/top-streams", s.topStreams)
	mux.HandleFunc("GET /api/stats/overview", s.statsOverview)
	mux.HandleFunc("GET /api/stats/search", s.searchLogs)
}

func (s *Server) statsByStream(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.StatsByStream()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) statsByLevel(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.StatsByLevel()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) topStreams(w http.ResponseWriter, r *http.Request) {
	n, _ := strconv.Atoi(r.URL.Query().Get("n"))
	if n < 1 {
		n = 5
	}
	result, err := s.svc.TopStreams(n)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) statsOverview(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.Overview()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) searchLogs(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	query := r.URL.Query().Get("q")
	var startTime, endTime time.Time
	if v := r.URL.Query().Get("start"); v != "" {
		startTime, _ = time.Parse(time.RFC3339, v)
	}
	if v := r.URL.Query().Get("end"); v != "" {
		endTime, _ = time.Parse(time.RFC3339, v)
	}
	items, total, err := s.svc.SearchLogs(query, startTime, endTime, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}
