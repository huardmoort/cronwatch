package heartbeat_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/example/cronwatch/internal/heartbeat"
	"github.com/example/cronwatch/internal/store"
)

func tempStore(t *testing.T) *store.Store {
	t.Helper()
	p := filepath.Join(t.TempDir(), "state.json")
	s, err := store.New(p)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	t.Cleanup(func() { os.Remove(p) })
	return s
}

func TestHandler_MissingJobParam(t *testing.T) {
	r := heartbeat.New(tempStore(t))
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	r.Handler()(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandler_RecordsHeartbeat(t *testing.T) {
	s := tempStore(t)
	r := heartbeat.New(s)

	before := time.Now().Add(-time.Second)
	req := httptest.NewRequest(http.MethodGet, "/ping?job=my-job", nil)
	w := httptest.NewRecorder()
	r.Handler()(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "my-job") {
		t.Errorf("response body missing job name: %s", w.Body.String())
	}

	last, ok := s.LastRun("my-job")
	if !ok {
		t.Fatal("expected LastRun to be recorded")
	}
	if last.Before(before) {
		t.Errorf("recorded time %v is before test start %v", last, before)
	}
}

func TestHandler_MultipleJobs(t *testing.T) {
	s := tempStore(t)
	r := heartbeat.New(s)
	handler := r.Handler()

	for _, job := range []string{"job-a", "job-b"} {
		req := httptest.NewRequest(http.MethodGet, "/ping?job="+job, nil)
		w := httptest.NewRecorder()
		handler(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("job %s: expected 200, got %d", job, w.Code)
		}
	}

	for _, job := range []string{"job-a", "job-b"} {
		if _, ok := s.LastRun(job); !ok {
			t.Errorf("expected LastRun recorded for %s", job)
		}
	}
}
