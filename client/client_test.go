package client

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFetchFile_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("file con"))
	}))
	defer server.Close()

	client := NewClient(server.Client())

	reader, err := client.FetchFile(server.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	content := &bytes.Buffer{}

	_, err = io.Copy(content, reader)
	if err != nil {
		t.Fatalf("faild to read content: %v", err)
	}

	if content.String() != "file con" {
		t.Errorf("Expected 'file con', got '%s'", content.String())
	}
}

func TestFetchFile_NetworkError(t *testing.T) {
	client := NewClient(&http.Client{})

	_, err := client.FetchFile("http://invalid.url")
	if err == nil {
		t.Fatal("Expected error for invalid URL")
	}

	if !strings.Contains(err.Error(), "request failed") {
		t.Errorf("Expected error to contain 'request failed', got %v", err)
	}
}

func TestFetchFile_Status404(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(server.Client())

	_, err := client.FetchFile(server.URL)
	if err == nil {
		t.Fatal("Expected error for 404 status")
	}

	expected := "file not found (404)"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("Expected error to contain '%s', got %v", expected, err)
	}
}

func TestFetchFile_Status403(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	client := NewClient(server.Client())

	_, err := client.FetchFile(server.URL)
	if err == nil {
		t.Fatal("Expected error for 403 status")
	}

	expected := "access forbidden (403)"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("Expected error to contain '%s', got %v", expected, err)
	}
}

func TestFetchFile_Status401(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewClient(server.Client())

	_, err := client.FetchFile(server.URL)
	if err == nil {
		t.Fatal("Expected error for 401 status")
	}

	expected := "unauthorized (401)"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("Expected error to contain '%s', got %v", expected, err)
	}
}

func TestFetchFile_Status429(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	client := NewClient(server.Client())

	_, err := client.FetchFile(server.URL)
	if err == nil {
		t.Fatal("Expected error for 429 status")
	}

	expected := "rate limited (429)"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("Expected error to contain '%s', got %v", expected, err)
	}
}

func TestFetchFile_OtherStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(server.Client())

	_, err := client.FetchFile(server.URL)
	if err == nil {
		t.Fatal("Expected error for 500 status")
	}

	expected := "request failed with status"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("Expected error to contain '%s', got %v", expected, err)
	}
}
