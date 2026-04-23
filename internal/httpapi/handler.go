package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/eelcoh/activities/internal/activity"
)

type Store interface {
	List(context.Context) ([]activity.Activity, error)
	Get(context.Context, string) (activity.Activity, error)
	Create(context.Context, activity.Activity) error
	Update(context.Context, activity.Activity) error
}

type Handler struct {
	store   Store
	service activity.Service
}

func NewHandler(store Store, service activity.Service) *Handler {
	return &Handler{
		store:   store,
		service: service,
	}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("GET /activities", h.listActivities)
	mux.HandleFunc("GET /activities/comments/{actUUID}", h.getComment)
	mux.HandleFunc("GET /activities/blogs/{actUUID}", h.getBlog)
	mux.HandleFunc("POST /activities/comments", h.createComment)
	mux.HandleFunc("POST /activities/blogs", h.createBlog)
	mux.HandleFunc("PUT /activities/comments/{actUUID}", h.updateComment)
	mux.HandleFunc("PUT /activities/blogs/{actUUID}", h.updateBlog)

	return mux
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Up!\n")
}

func (h *Handler) listActivities(w http.ResponseWriter, r *http.Request) {
	activities, err := h.store.List(r.Context())
	if err != nil {
		http.Error(w, "activities not found", http.StatusServiceUnavailable)
		return
	}

	writeJSON(w, http.StatusOK, activity.ListResponse{Activities: activities})
}

func (h *Handler) getComment(w http.ResponseWriter, r *http.Request) {
	h.getActivity(w, r, activity.TypeComment)
}

func (h *Handler) getBlog(w http.ResponseWriter, r *http.Request) {
	h.getActivity(w, r, activity.TypeBlog)
}

func (h *Handler) getActivity(w http.ResponseWriter, r *http.Request, wantType string) {
	item, err := h.store.Get(r.Context(), r.PathValue("actUUID"))
	if errors.Is(err, activity.ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	if item.Type() != wantType {
		http.Error(w, "incorrect activity type", http.StatusForbidden)
		return
	}

	writeJSON(w, http.StatusOK, item)
}

func (h *Handler) createComment(w http.ResponseWriter, r *http.Request) {
	var payload activity.CommentEnvelope
	if err := decodeJSON(r, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	comment, err := h.service.NewComment(payload.Comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.store.Create(r.Context(), comment); err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	writeJSON(w, http.StatusCreated, comment)
}

func (h *Handler) createBlog(w http.ResponseWriter, r *http.Request) {
	var payload activity.BlogEnvelope
	if err := decodeJSON(r, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	blog, err := h.service.NewBlog(payload.Blog)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.store.Create(r.Context(), blog); err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	writeJSON(w, http.StatusCreated, blog)
}

func (h *Handler) updateComment(w http.ResponseWriter, r *http.Request) {
	var input activity.CommentInput
	if err := decodeJSON(r, &input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item, err := h.store.Get(r.Context(), r.PathValue("actUUID"))
	if errors.Is(err, activity.ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	comment, ok := item.(activity.Comment)
	if !ok {
		http.Error(w, "incorrect activity type", http.StatusForbidden)
		return
	}

	updated, err := h.service.UpdateComment(comment, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.store.Update(r.Context(), updated); err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (h *Handler) updateBlog(w http.ResponseWriter, r *http.Request) {
	var input activity.BlogInput
	if err := decodeJSON(r, &input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item, err := h.store.Get(r.Context(), r.PathValue("actUUID"))
	if errors.Is(err, activity.ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	blog, ok := item.(activity.Blog)
	if !ok {
		http.Error(w, "incorrect activity type", http.StatusForbidden)
		return
	}

	updated, err := h.service.UpdateBlog(blog, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.store.Update(r.Context(), updated); err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func decodeJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return errors.New("request body is required")
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
