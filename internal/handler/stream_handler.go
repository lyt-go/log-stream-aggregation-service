package handler

import (
	"net/http"

	"logaggregation/internal/model"
	"logaggregation/pkg/httpx"
)

func (s *Server) registerStreamRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/streams", s.createStream)
	mux.HandleFunc("GET /api/streams", s.listStreams)
	mux.HandleFunc("GET /api/streams/{id}", s.getStream)
	mux.HandleFunc("PUT /api/streams/{id}", s.updateStream)
	mux.HandleFunc("DELETE /api/streams/{id}", s.deleteStream)
	mux.HandleFunc("POST /api/streams/{id}/pause", s.pauseStream)
	mux.HandleFunc("POST /api/streams/{id}/resume", s.resumeStream)
}

type createStreamRequest struct {
	Name   string `json:"name"`
	Source string `json:"source"`
	Format string `json:"format"`
	Status string `json:"status"`
}

func (s *Server) createStream(w http.ResponseWriter, r *http.Request) {
	var req createStreamRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	st, err := s.svc.CreateStream(model.Stream{Name: req.Name, Source: req.Source, Format: req.Format, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, st)
}

func (s *Server) listStreams(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.StreamFilter{
		Status:  r.URL.Query().Get("status"),
		Format:  r.URL.Query().Get("format"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListStreams(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getStream(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	st, err := s.svc.GetStream(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, st)
}

type updateStreamRequest struct {
	Name   string `json:"name"`
	Source string `json:"source"`
	Format string `json:"format"`
}

func (s *Server) updateStream(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateStreamRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	st, err := s.svc.UpdateStream(id, model.Stream{Name: req.Name, Source: req.Source, Format: req.Format})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, st)
}

func (s *Server) deleteStream(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteStream(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) pauseStream(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	st, err := s.svc.TransitionStreamStatus(id, model.StreamStatusPaused)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, st)
}

func (s *Server) resumeStream(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	st, err := s.svc.TransitionStreamStatus(id, model.StreamStatusActive)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, st)
}
