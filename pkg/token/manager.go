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
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"io"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

var ErrCipherTextTooShort = errors.New("ciphertext too short")

type TokenManager struct {
	gcm cipher.AEAD
}

func NewTokenManager(key []byte) *TokenManager {
	// Padding key to 32 bytes
	if len(key) < 32 {
		padding := make([]byte, 32-len(key))
		key = append(key, padding...)
	}
	if len(key) > 32 {
		logrus.Warn("Key length is longer than 32 bytes, will be truncated")
		key = key[:32]
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		panic(err)
	}

	return &TokenManager{gcm: gcm}
}

func (tm *TokenManager) encrypt(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, tm.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return tm.gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func (tm *TokenManager) decrypt(ciphertext []byte) ([]byte, error) {
	nonceSize := tm.gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, ErrCipherTextTooShort
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return tm.gcm.Open(nil, nonce, ciphertext, nil)
}

func (tm *TokenManager) Encrypt(token *AgentToken) (string, error) {
	// serialize token
	data, err := proto.Marshal(token)
	if err != nil {
		return "", err
	}
	// encrypt
	encData, err := tm.encrypt(data)
	if err != nil {
		return "", err
	}

	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(encData), nil
}

func (tm *TokenManager) Decrypt(encrypted string) (*AgentToken, error) {
	encData, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(encrypted)
	if err != nil {
		return nil, err
	}

	// decrypt
	data, err := tm.decrypt(encData)
	if err != nil {
		return nil, err
	}
	// deserialize token
	token := &AgentToken{}
	err = proto.Unmarshal(data, token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

var agentTokenKey int

func NewContext(ctx context.Context, token *AgentToken) context.Context {
	return context.WithValue(ctx, &agentTokenKey, token)
}

func FromContext(ctx context.Context) *AgentToken {
	token, _ := ctx.Value(&agentTokenKey).(*AgentToken)
	return token
}
