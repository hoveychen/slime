package hub

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hoveychen/slime/pkg/token"
)

type mockTokenManager struct {
	tok *token.AgentToken
}

func (m *mockTokenManager) Encrypt(t *token.AgentToken) (string, error) {
	return "encrypted-token", nil
}

func (m *mockTokenManager) Decrypt(s string) (*token.AgentToken, error) {
	if s != "encrypted-token" {
		return nil, errors.New("invalid token")
	}
	return m.tok, nil
}

func TestWrapTokenValidator(t *testing.T) {
	// Create a mock token manager
	tokenMgr := &mockTokenManager{}

	// Create a mock handler that just writes a response
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create a test server with the wrapTokenValidator middleware
	testServer := httptest.NewServer((&HubServer{tokenMgr: tokenMgr}).wrapTokenValidator(mockHandler))
	defer testServer.Close()

	// Test case 1: valid tok
	tok := &token.AgentToken{ExpireAt: time.Now().Add(time.Hour).Unix()}
	tokenMgr.tok = tok
	encryptedToken, err := tokenMgr.Encrypt(tok)
	if err != nil {
		t.Fatalf("Failed to encrypt token: %v", err)
	}
	req, err := http.NewRequest("GET", testServer.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("slime-agent-token", encryptedToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Test case 2: invalid token
	encryptedToken = "invalid-token"
	req, err = http.NewRequest("GET", testServer.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("slime-agent-token", encryptedToken)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}

	// Test case 3: expired token
	tok = &token.AgentToken{ExpireAt: time.Now().Add(-time.Hour).Unix()}
	tokenMgr.tok = tok
	encryptedToken, err = tokenMgr.Encrypt(tok)
	if err != nil {
		t.Fatalf("Failed to encrypt token: %v", err)
	}
	req, err = http.NewRequest("GET", testServer.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("slime-agent-token", encryptedToken)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}
