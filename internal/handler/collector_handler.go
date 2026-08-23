package handler

import (
	"net/http"

	"logaggregation/internal/model"
	"logaggregation/pkg/httpx"
)

func (s *Server) registerCollectorRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/collectors", s.createCollector)
	mux.HandleFunc("GET /api/collectors", s.listCollectors)
	mux.HandleFunc("GET /api/collectors/{id}", s.getCollector)
	mux.HandleFunc("PUT /api/collectors/{id}", s.updateCollector)
	mux.HandleFunc("DELETE /api/collectors/{id}", s.deleteCollector)
	mux.HandleFunc("POST /api/collectors/{id}/stop", s.stopCollector)
	mux.HandleFunc("POST /api/collectors/{id}/start", s.startCollector)
}

type createCollectorRequest struct {
	Name     string `json:"name"`
	StreamID string `json:"stream_id"`
	Endpoint string `json:"endpoint"`
	Status   string `json:"status"`
}

func (s *Server) createCollector(w http.ResponseWriter, r *http.Request) {
	var req createCollectorRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.CreateCollector(model.Collector{Name: req.Name, StreamID: req.StreamID, Endpoint: req.Endpoint, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, c)
}

func (s *Server) listCollectors(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.CollectorFilter{
		Status:   r.URL.Query().Get("status"),
		StreamID: r.URL.Query().Get("stream_id"),
		Keyword:  r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListCollectors(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getCollector(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c, err := s.svc.GetCollector(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

type updateCollectorRequest struct {
	Name     string `json:"name"`
	StreamID string `json:"stream_id"`
	Endpoint string `json:"endpoint"`
}

func (s *Server) updateCollector(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateCollectorRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.UpdateCollector(id, model.Collector{Name: req.Name, StreamID: req.StreamID, Endpoint: req.Endpoint})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) deleteCollector(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteCollector(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) stopCollector(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c, err := s.svc.TransitionCollectorStatus(id, model.CollectorStatusStopped)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) startCollector(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c, err := s.svc.TransitionCollectorStatus(id, model.CollectorStatusActive)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}
