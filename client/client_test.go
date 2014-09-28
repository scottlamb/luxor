package client_test

import (
	"code.google.com/p/go.net/context"
	"github.com/scottlamb/luxor/client"
	"github.com/scottlamb/luxor/protocol"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestSuccessfulThemeGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ThemeGet.json" {
			t.Errorf("expected /ThemeGet.json; got %v", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %v", r.Method)
		}
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected json; got %v", contentType)
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Body read failed; %v", err)
		}
		expectedBody := "{\"ThemeIndex\":0}"
		stringBody := string(body)
		if stringBody != expectedBody {
			t.Errorf("Expected %v; got %v", expectedBody, stringBody)
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w,
			"{\"Groups\":["+
				"{\"GroupNumber\":0,\"Intensity\":26},"+
				"{\"GroupNumber\":1,\"Intensity\":42}"+
				"]}")
	}))
	defer server.Close()
	c := &client.Controller{BaseURL: server.URL}
	req := &protocol.ThemeGetRequest{ThemeIndex: 0}
	resp, err := c.ThemeGet(context.Background(), req)
	if err != nil {
		t.Errorf("expected success; got %v", err)
		return
	}
	if len(resp.Groups) != 2 {
		t.Errorf("Expected 2 groups; got %+v", resp.Groups)
		return
	}
	if resp.Groups[0].GroupNumber != 0 || resp.Groups[0].Intensity != 26 {
		t.Errorf("Group 0 mismatch; groups are %+v", resp.Groups)
	}
	if resp.Groups[1].GroupNumber != 1 || resp.Groups[1].Intensity != 42 {
		t.Errorf("Group 1 mismatch; groups are %+v", resp.Groups)
	}
}

func TestAlreadyCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c := &client.Controller{BaseURL: ""}
	_, err := c.GroupListGet(ctx, &protocol.GroupListGetRequest{})
	if err != context.Canceled {
		t.Errorf("expected canceled; got %v", err)
	}
}

// A context which has a Deadline but does not become Done when it passes.
type deadlineContext time.Time

func (c deadlineContext) Deadline() (deadline time.Time, ok bool) {
	return time.Time(c), true
}

func (c deadlineContext) Done() <-chan struct{} {
	return make(chan struct{})
}

func (c deadlineContext) Err() error {
	return nil
}

func (c deadlineContext) Value(key interface{}) interface{} {
	return nil
}

func TestHttpTimeoutLe0(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()
	c := &client.Controller{BaseURL: server.URL}
	req := &protocol.ThemeGetRequest{ThemeIndex: 0}
	_, err := c.ThemeGet(deadlineContext(time.Time{}), req)
	if err != context.Canceled {
		t.Errorf("expected canceled; got %v", err)
	}
}

func TestHttpTimeoutGt0(t *testing.T) {
	requestDone := make(chan struct{}, 1) // closed when request finishes.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-requestDone:
		case <-time.After(time.Minute):
			t.Error("took way too long to finish")
		}
	}))
	defer server.Close()
	c := &client.Controller{BaseURL: server.URL}
	req := &protocol.ThemeGetRequest{ThemeIndex: 0}
	_, err := c.ThemeGet(deadlineContext(time.Now().Add(time.Second)), req)
	close(requestDone)
	if err == nil {
		t.Error("should have failed")
	}
}

func TestCanceledWhileWaiting(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	requestDone := make(chan struct{}, 1) // closed when request finishes.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cancel()
		select {
		case <-requestDone:
		case <-time.After(time.Minute):
			t.Error("took way too long to finish")
		}
	}))
	defer server.Close()
	c := &client.Controller{BaseURL: server.URL}
	req := &protocol.ThemeGetRequest{ThemeIndex: 0}
	_, err := c.ThemeGet(ctx, req)
	close(requestDone)
	if err != context.Canceled {
		t.Errorf("expected canceled; got %v", err)
	}
}

func TestPostFailed(t *testing.T) {
	c := &client.Controller{BaseURL: "badscheme://"}
	req := &protocol.ThemeGetRequest{ThemeIndex: 0}
	_, err := c.ThemeGet(context.Background(), req)
	if !strings.Contains(err.Error(), "badscheme") {
		t.Errorf("expected http error about scheme, got %v", err)
	}
}

func TestBadStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Horrible problem", http.StatusInternalServerError)
	}))
	defer server.Close()
	c := &client.Controller{BaseURL: server.URL}
	req := &protocol.ThemeGetRequest{ThemeIndex: 0}
	resp, err := c.ThemeGet(context.Background(), req)
	if err == nil {
		t.Errorf("expected error; got success: %+v", resp)
		return
	}
	if !strings.Contains(err.Error(), "Horrible problem") {
		t.Errorf("Error should mention horrible problem; %v", err)
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("Error should mention HTTP status code; %v", err)
	}
}

func TestBadContentType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w,
			"{\"Groups\":["+
				"{\"GroupNumber\":0,\"Intensity\":26},"+
				"{\"GroupNumber\":1,\"Intensity\":42}"+
				"]}")
	}))
	defer server.Close()
	c := &client.Controller{BaseURL: server.URL}
	req := &protocol.ThemeGetRequest{ThemeIndex: 0}
	resp, err := c.ThemeGet(context.Background(), req)
	if err == nil {
		t.Errorf("expected error; got success: %+v", resp)
		return
	}
	if !strings.Contains(err.Error(), "text/plain") {
		t.Errorf("Error should mention MIME type; %v", err)
	}
}

func TestBadJson(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, "asdf")
	}))
	defer server.Close()
	c := &client.Controller{BaseURL: server.URL}
	req := &protocol.ThemeGetRequest{ThemeIndex: 0}
	resp, err := c.ThemeGet(context.Background(), req)
	if err == nil {
		t.Errorf("expected error; got success: %+v", resp)
		return
	}
	if !strings.Contains(err.Error(), "asdf") {
		t.Errorf("Error should include bogus body; %v", err)
	}
}

func TestBadApplicationStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, "{\"Status\":1}")
	}))
	defer server.Close()
	c := &client.Controller{BaseURL: server.URL}
	req := &protocol.ThemeGetRequest{ThemeIndex: 0}
	resp, err := c.ThemeGet(context.Background(), req)
	if err == nil {
		t.Errorf("expected error; got success: %+v", resp)
		return
	}
	if !strings.Contains(err.Error(), "unknown method") {
		t.Errorf("Error should describe status; %v", err)
	}
}
