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
package token

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenManager_EncryptDecrypt(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef")
	tm := NewTokenManager(key)

	plaintext := []byte("hello, world!")
	ciphertext, err := tm.encrypt(plaintext)
	assert.NoError(t, err)

	decrypted, err := tm.decrypt(ciphertext)
	assert.NoError(t, err)
	assert.True(t, bytes.Equal(plaintext, decrypted))
}

func TestTokenManager_EncryptDecrypt_Empty(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef")
	tm := NewTokenManager(key)

	plaintext := []byte("")
	ciphertext, err := tm.encrypt(plaintext)
	assert.NoError(t, err)

	decrypted, err := tm.decrypt(ciphertext)
	assert.NoError(t, err)
	assert.True(t, bytes.Equal(plaintext, decrypted))
}

func TestTokenManager_EncryptDecrypt_TooShort(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef")
	tm := NewTokenManager(key)

	ciphertext := []byte("too short")
	_, err := tm.decrypt(ciphertext)
	assert.Equal(t, ErrCipherTextTooShort, err)
}

func TestTokenManager_EncryptDecrypt_WrongKey(t *testing.T) {
	key1 := []byte("0123456789abcdef0123456789abcdef")
	key2 := []byte("0123456789abcdef0123456789abcdee")
	tm1 := NewTokenManager(key1)
	tm2 := NewTokenManager(key2)

	plaintext := []byte("hello, world!")
	ciphertext, err := tm1.encrypt(plaintext)
	assert.NoError(t, err)

	_, err = tm2.decrypt(ciphertext)
	assert.Error(t, err)
}

func TestTokenManager_EncryptDecryptToken(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef")
	tm := NewTokenManager(key)

	token := &AgentToken{
		Id:         123,
		Name:       "test",
		ScopePaths: []string{"/api/foo/bar"},
	}

	encrypted, err := tm.Encrypt(token)
	assert.NoError(t, err)

	decrypted, err := tm.Decrypt(encrypted)
	assert.NoError(t, err)
	if token.GetId() != decrypted.GetId() {
		t.Errorf("id not match: %d != %d", token.GetId(), decrypted.GetId())
	}
	if token.GetName() != decrypted.GetName() {
		t.Errorf("name not match: %s != %s", token.GetName(), decrypted.GetName())
	}
	if len(token.GetScopePaths()) != len(decrypted.GetScopePaths()) {
		t.Errorf("scope paths not match: %v != %v", token.GetScopePaths(), decrypted.GetScopePaths())
	}
	for i := range token.GetScopePaths() {
		if token.GetScopePaths()[i] != decrypted.GetScopePaths()[i] {
			t.Errorf("scope paths not match: %v != %v", token.GetScopePaths(), decrypted.GetScopePaths())
		}
	}
}

func TestNewContext(t *testing.T) {
	// Create a new context
	ctx := context.Background()

	// Create a new agent token
	token := &AgentToken{
		Id:         123,
		Name:       "test",
		ScopePaths: []string{"/api/foo/bar"},
	}

	// Add the agent token to the context
	ctx = NewContext(ctx, token)

	// Retrieve the agent token from the context
	retrievedToken := FromContext(ctx)

	// Check that the retrieved token matches the original token
	assert.Equal(t, token, retrievedToken)
}

func TestFromContext_NoToken(t *testing.T) {
	// Create a new context
	ctx := context.Background()

	// Retrieve the agent token from the context
	retrievedToken := FromContext(ctx)

	// Check that the retrieved token is nil
	assert.Nil(t, retrievedToken)
}

func TestNewTokenManager(t *testing.T) {
	// Test with key length less than 32 bytes
	key := []byte("0123456789abcdef")
	tm := NewTokenManager(key)
	assert.NotNil(t, tm.gcm)

	// Test with key length equal to 32 bytes
	key = []byte("0123456789abcdef0123456789abcdef")
	tm = NewTokenManager(key)
	assert.NotNil(t, tm.gcm)

	// Test with key length greater than 32 bytes
	key = []byte("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0")
	tm = NewTokenManager(key)
	assert.NotNil(t, tm.gcm)
}
