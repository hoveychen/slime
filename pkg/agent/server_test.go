/*
Copyright © 2023 Harry C <hoveychen@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package agent

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
	req := as.newHubAPIRequest(context.Background(), 123, apiPath, reader)

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

func TestNewAgentServer(t *testing.T) {
	// Test case 1: Valid hub and upstream addresses, with default options.
	hubAddr := "http://localhost:8080"
	upstreamAddr := "http://localhost:9090"
	token := "abc123"
	as, err := NewAgentServer(hubAddr, upstreamAddr, token)
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}
	if as.token != token {
		t.Errorf("Expected token to be %s, but got %s", token, as.token)
	}
	if as.numWorker != defaultNumWorker {
		t.Errorf("Expected numWorker to be %d, but got %d", defaultNumWorker, as.numWorker)
	}
	if as.hubURL.String() != hubAddr {
		t.Errorf("Expected hubURL to be %s, but got %s", hubAddr, as.hubURL.String())
	}
	if as.upstreamURL.String() != upstreamAddr {
		t.Errorf("Expected upstreamURL to be %s, but got %s", upstreamAddr, as.upstreamURL.String())
	}

	// Test case 2: Valid hub and upstream addresses, with custom options.
	numWorker := 5
	as, err = NewAgentServer(hubAddr, upstreamAddr, token, WithNumWorker(numWorker))
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}
	if as.numWorker != numWorker {
		t.Errorf("Expected numWorker to be %d, but got %d", numWorker, as.numWorker)
	}

	// Test case 3: Invalid hub address.
	hubAddr = "invalid"
	_, err = NewAgentServer(hubAddr, upstreamAddr, token)
	if err == nil {
		t.Errorf("Expected non-nil error, but got nil")
	}

	// Test case 4: Invalid upstream address.
	hubAddr = "http://localhost:8080"
	upstreamAddr = "invalid"
	_, err = NewAgentServer(hubAddr, upstreamAddr, token)
	if err == nil {
		t.Errorf("Expected non-nil error, but got nil")
	}
}

func TestWithNumWorker(t *testing.T) {
	// Create a new AgentServer with default options.
	hubAddr := "http://localhost:8080"
	upstreamAddr := "http://localhost:9090"
	token := "abc123"
	as, err := NewAgentServer(hubAddr, upstreamAddr, token)
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}

	// Call the WithNumWorker option with a custom number of workers.
	numWorker := 5
	WithNumWorker(numWorker)(as)
	if as.numWorker != numWorker {
		t.Errorf("Expected numWorker to be %d, but got %d", numWorker, as.numWorker)
	}
}

func TestWithReportHardware(t *testing.T) {
	// Create a new AgentServer with reportHW set to false.
	as := &AgentServer{}

	// Call the WithReportHardware function with true.
	WithReportHardware(true)(as)

	// Check that reportHW is true.
	if !as.reportHW {
		t.Errorf("Expected reportHW to be true, but got false")
	}

	// Call the WithReportHardware function with false.
	WithReportHardware(false)(as)

	// Check that reportHW is false.
	if as.reportHW {
		t.Errorf("Expected reportHW to be false, but got true")
	}
}

func TestJoinHubSuccess(t *testing.T) {
	// Create a new AgentServer with a mock hub URL and token.
	as := &AgentServer{
		hubURL: &url.URL{Scheme: "http", Host: "localhost:8080"},
		token:  "abc123",
	}

	// Create a new mock HTTP server that returns a 200 OK response.
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	// Set the mock server URL as the hub URL.
	as.hubURL, _ = url.Parse(mockServer.URL)

	// Call the joinHub method.
	err := as.joinHub(context.Background(), 123)

	// Check that the error is nil.
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}
}

func TestJoinHubError(t *testing.T) {
	// Create a new AgentServer with a mock hub URL and token.
	as := &AgentServer{
		hubURL: &url.URL{Scheme: "http", Host: "localhost:8080"},
		token:  "abc123",
	}

	// Create a new mock HTTP server that returns a 500 Internal Server Error response.
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	// Set the mock server URL as the hub URL.
	as.hubURL, _ = url.Parse(mockServer.URL)

	// Call the joinHub method.
	err := as.joinHub(context.Background(), 123)

	// Check that the error is not nil.
	if err == nil {
		t.Errorf("Expected non-nil error, but got nil")
	}
}

func TestJoinHubHTTPError(t *testing.T) {
	// Create a new AgentServer with a mock hub URL and token.
	as := &AgentServer{
		hubURL: &url.URL{Scheme: "http", Host: "localhost:8080"},
		token:  "abc123",
	}

	// Create a new mock HTTP server that returns a 404 Not Found response.
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	// Set the mock server URL as the hub URL.
	as.hubURL, _ = url.Parse(mockServer.URL)

	// Call the joinHub method.
	err := as.joinHub(context.Background(), 123)

	// Check that the error is not nil.
	if err == nil {
		t.Errorf("Expected non-nil error, but got nil")
	}
}

func TestGetOrCreateAgentID(t *testing.T) {
	// Backup the original agentIDFile and restore it after the test
	originalAgentIDFile := agentIDFile
	defer func() { agentIDFile = originalAgentIDFile }()

	// Create a temporary file for testing
	tmpfile, err := os.CreateTemp("", "agentID")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	tmpfile.Close()
	os.Remove(tmpfile.Name())

	agentIDFile = tmpfile.Name()

	// Test case 1: agent ID file does not exist
	id, err := getOrCreateAgentID()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	data, err := os.ReadFile(agentIDFile)
	if err != nil {
		t.Fatalf("Failed to read agent ID file: %v", err)
	}
	expectedID, err := strconv.Atoi(string(data))
	if err != nil {
		t.Fatalf("Invalid agent ID in file: %v", err)
	}
	if id != expectedID {
		t.Errorf("Expected ID %d, got %d", expectedID, id)
	}

	// Test case 2: agent ID file exists
	id, err = getOrCreateAgentID()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if id != expectedID {
		t.Errorf("Expected ID %d, got %d", expectedID, id)
	}

	// Test case 3: agent ID file contains invalid data
	err = os.WriteFile(agentIDFile, []byte("invalid"), 0644)
	if err != nil {
		t.Fatalf("Failed to write to agent ID file: %v", err)
	}
	_, err = getOrCreateAgentID()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestWithAgentID(t *testing.T) {
	// Create a mock AgentServer
	as := &AgentServer{}

	// Call the WithAgentID function
	WithAgentID(123)(as)

	// Check that the agentID is set correctly
	assert.Equal(t, 123, as.agentID)
}
