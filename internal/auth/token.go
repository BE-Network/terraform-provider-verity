package auth

import (
	"net/http"
	"sync"
	"time"
)

type TokenManager struct {
	mu            sync.RWMutex
	token         string
	expiresAt     time.Time
	cookieJar     http.CookieJar
	refreshBefore time.Duration
}

func NewTokenManager(jar http.CookieJar) *TokenManager {
	return &TokenManager{
		cookieJar:     jar,
		refreshBefore: 5 * time.Minute,
	}
}

func (tm *TokenManager) SetToken(token string, expiresIn time.Duration) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.token = token
	tm.expiresAt = time.Now().Add(expiresIn)
}

func (tm *TokenManager) GetToken() (string, bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	if tm.token == "" {
		return "", true
	}

	needsRefresh := time.Now().After(tm.expiresAt)
	return tm.token, needsRefresh
}

func (tm *TokenManager) GetCookieJar() http.CookieJar {
	return tm.cookieJar
}

func (tm *TokenManager) Clear() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.token = ""
	tm.expiresAt = time.Time{}
}
