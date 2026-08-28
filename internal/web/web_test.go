package web

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"qin-culture-site/internal/catalog"
	"qin-culture-site/internal/service"
	"qin-culture-site/internal/store"
)

func TestWebHomeAndExperience(t *testing.T) {
	db, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	srv, err := NewServer(service.New(catalog.New(), db), 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	home := httptest.NewRecorder()
	srv.Handler().ServeHTTP(home, httptest.NewRequest(http.MethodGet, "/", nil))
	if home.Code != http.StatusOK || !strings.Contains(home.Body.String(), "古琴艺术专题") || !strings.Contains(home.Body.String(), "图片墙") {
		t.Fatalf("unexpected home response: %d", home.Code)
	}
	request := httptest.NewRequest(http.MethodPost, "/api/experience", strings.NewReader(`{"name":"访客","contact":"mail","interest":"琴派故事","message":"请联系"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	srv.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("unexpected create status: %d %s", response.Code, response.Body.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil || payload["status"] != "pending" {
		t.Fatalf("unexpected payload: %#v", payload)
	}
	_ = context.Background()
}
