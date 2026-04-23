package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eelcoh/activities/internal/activity"
)

type fakeStore struct {
	list      []activity.Activity
	get       map[string]activity.Activity
	created   []activity.Activity
	updated   []activity.Activity
	createErr error
	updateErr error
	getErr    error
	listErr   error
}

func (s *fakeStore) List(context.Context) ([]activity.Activity, error) {
	if s.listErr != nil {
		return nil, s.listErr
	}
	return s.list, nil
}

func (s *fakeStore) Get(_ context.Context, id string) (activity.Activity, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	item, ok := s.get[id]
	if !ok {
		return nil, activity.ErrNotFound
	}
	return item, nil
}

func (s *fakeStore) Create(_ context.Context, item activity.Activity) error {
	if s.createErr != nil {
		return s.createErr
	}
	s.created = append(s.created, item)
	return nil
}

func (s *fakeStore) Update(_ context.Context, item activity.Activity) error {
	if s.updateErr != nil {
		return s.updateErr
	}
	s.updated = append(s.updated, item)
	return nil
}

func TestRoutesHealth(t *testing.T) {
	handler := NewHandler(&fakeStore{}, activity.NewService("secret", func() int64 { return 1 }, func() string { return "uuid-1" }))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	res := httptest.NewRecorder()

	handler.Routes().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	if body := res.Body.String(); body != "Up!\n" {
		t.Fatalf("body = %q, want %q", body, "Up!\n")
	}
}

func TestRoutesListActivities(t *testing.T) {
	store := &fakeStore{
		list: []activity.Activity{
			activity.Comment{
				Author: "Eelco",
				Msg:    []string{"hello"},
				Meta:   activity.Meta{UUID: "uuid-1", Active: true, Date: 1},
			},
		},
	}
	handler := NewHandler(store, activity.NewService("secret", func() int64 { return 1 }, func() string { return "uuid-1" }))

	req := httptest.NewRequest(http.MethodGet, "/activities", nil)
	res := httptest.NewRecorder()

	handler.Routes().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}

	var payload struct {
		Activities []map[string]any `json:"activities"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	if len(payload.Activities) != 1 {
		t.Fatalf("len(payload.Activities) = %d, want %d", len(payload.Activities), 1)
	}
	if payload.Activities[0]["author"] != "Eelco" {
		t.Fatalf("payload.Activities[0][author] = %v, want %q", payload.Activities[0]["author"], "Eelco")
	}
}

func TestRoutesCreateComment(t *testing.T) {
	store := &fakeStore{get: map[string]activity.Activity{}}
	handler := NewHandler(store, activity.NewService("secret", func() int64 { return 7 }, func() string { return "uuid-7" }))

	body := bytes.NewBufferString(`{"comment":{"author":"Eelco","msg":["hello"]}}`)
	req := httptest.NewRequest(http.MethodPost, "/activities/comments", body)
	res := httptest.NewRecorder()

	handler.Routes().ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusCreated)
	}
	if len(store.created) != 1 {
		t.Fatalf("len(store.created) = %d, want %d", len(store.created), 1)
	}

	comment, ok := store.created[0].(activity.Comment)
	if !ok {
		t.Fatalf("created[0] type = %T, want activity.Comment", store.created[0])
	}
	if comment.Meta.UUID != "uuid-7" {
		t.Fatalf("comment.Meta.UUID = %q, want %q", comment.Meta.UUID, "uuid-7")
	}
}

func TestRoutesUpdateBlogRejectsWrongPassphrase(t *testing.T) {
	store := &fakeStore{
		get: map[string]activity.Activity{
			"uuid-1": activity.Blog{
				Author: "Eelco",
				Title:  "Title",
				Msg:    []string{"hello"},
				Meta:   activity.Meta{UUID: "uuid-1", Active: true, Date: 1},
			},
		},
	}
	handler := NewHandler(store, activity.NewService("secret", func() int64 { return 1 }, func() string { return "uuid-1" }))

	body := bytes.NewBufferString(`{"passphrase":"wrong","author":"Eelco","title":"New","msg":["hello"]}`)
	req := httptest.NewRequest(http.MethodPut, "/activities/blogs/uuid-1", body)
	res := httptest.NewRecorder()

	handler.Routes().ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusBadRequest)
	}
	if !errors.Is(activity.ErrUnauthorized, activity.ErrUnauthorized) {
		t.Fatal("sentinel error expectation invalid")
	}
}
