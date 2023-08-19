package agent

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestParseAddrWithHostAndPort(t *testing.T) {
	// Call the parseAddr function with a host and port.
	addr := "localhost:8080"
	u, err := parseAddr(addr)

	// Check that the error is nil.
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}

	// Check that the URL scheme is "http".
	if u.Scheme != "http" {
		t.Errorf("Expected scheme to be http, but got %s", u.Scheme)
	}

	// Check that the URL host is "localhost:8080".
	if u.Host != "localhost:8080" {
		t.Errorf("Expected host to be localhost:8080, but got %s", u.Host)
	}
}

func TestParseAddrWithSchemeHostAndPort(t *testing.T) {
	// Call the parseAddr function with a scheme, host, and port.
	addr := "https://localhost:8080"
	u, err := parseAddr(addr)

	// Check that the error is nil.
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}

	// Check that the URL scheme is "https".
	if u.Scheme != "https" {
		t.Errorf("Expected scheme to be https, but got %s", u.Scheme)
	}

	// Check that the URL host is "localhost:8080".
	if u.Host != "localhost:8080" {
		t.Errorf("Expected host to be localhost:8080, but got %s", u.Host)
	}
}

func TestParseAddrWithInvalidScheme(t *testing.T) {
	// Call the parseAddr function with an invalid scheme.
	addr := "ftp://localhost:8080"
	_, err := parseAddr(addr)

	// Check that the error is not nil.
	if err == nil {
		t.Errorf("Expected non-nil error, but got nil")
	}
}

func TestParseAddrWithInvalidHost(t *testing.T) {
	// Call the parseAddr function with an invalid host.
	addr := "http://"
	a, err := parseAddr(addr)
	fmt.Printf("%#v", a)

	// Check that the error is not nil.
	if err == nil {
		t.Errorf("Expected non-nil error, but got nil")
	}
}

func TestParseAddrWithInvalidAddr(t *testing.T) {
	// Call the parseAddr function with an invalid address.
	addr := "invalid"
	_, err := parseAddr(addr)

	// Check that the error is not nil.
	if err == nil {
		t.Errorf("Expected non-nil error, but got nil")
	}
}

func TestNewHubAPIRequest(t *testing.T) {
	// Create a new AgentServer with a mock hub URL and token.
	as := &AgentServer{
		hubURL: &url.URL{Scheme: "http", Host: "localhost:8080"},
		token:  "abc123",
	}

	// Create a new hub API request with a mock API path and reader.
	apiPath := "/api/v1/test"
	reader := strings.NewReader("test")
	req := as.newHubAPIRequest(context.Background(), apiPath, reader)

	// Check that the request method is "POST".
	if req.Method != "POST" {
		t.Errorf("Expected method to be POST, but got %s", req.Method)
	}

	// Check that the request URL has the expected path.
	if req.URL.Path != apiPath {
		t.Errorf("Expected path to be %s, but got %s", apiPath, req.URL.Path)
	}

	// Check that the request header has the expected token.
	if req.Header.Get("slime-agent-token") != as.token {
		t.Errorf("Expected token to be %s, but got %s", as.token, req.Header.Get("slime-agent-token"))
	}
}

func TestFixUpstreamRequest(t *testing.T) {
	// Create a new AgentServer with a mock upstream URL.
	as := &AgentServer{
		upstreamURL: &url.URL{Scheme: "http", Host: "localhost:8081", Path: "/api/v1"},
	}

	// Create a new upstream request with a mock path and host.
	path := "/test"
	host := "example.com"
	req := &http.Request{
		URL:  &url.URL{Scheme: "https", Host: host, Path: path},
		Host: host,
	}

	// Fix the upstream request.
	as.fixUpstreamRequest(req)

	// Check that the request URL has the expected scheme, host, and path.
	if req.URL.Scheme != as.upstreamURL.Scheme {
		t.Errorf("Expected scheme to be %s, but got %s", as.upstreamURL.Scheme, req.URL.Scheme)
	}
	if req.URL.Host != as.upstreamURL.Host {
		t.Errorf("Expected host to be %s, but got %s", as.upstreamURL.Host, req.URL.Host)
	}
	if req.URL.Path != path {
		t.Errorf("Expected path to be %s, but got %s", path, req.URL.Path)
	}

	// Check that the request host and request URI are empty.
	if req.Host != "" {
		t.Errorf("Expected host to be empty, but got %s", req.Host)
	}
	if req.RequestURI != "" {
		t.Errorf("Expected request URI to be empty, but got %s", req.RequestURI)
	}
}
