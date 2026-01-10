package rest

import (
	"net/http"
	"strings"
)

type TokenSet struct {
	enabled bool
	tokens  map[string]struct{} // value -> present
}

func NewTokenSet(enabled bool, tokenValues []string) *TokenSet {
	m := make(map[string]struct{}, len(tokenValues))
	for _, t := range tokenValues {
		if strings.TrimSpace(t) == "" {
			continue
		}
		m[t] = struct{}{}
	}
	return &TokenSet{enabled: enabled, tokens: m}
}

func (ts *TokenSet) Require(next http.Handler, stats *Stats) http.Handler {
	if !ts.enabled {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// /health is always public
		if r.URL.Path == "/api/v1/health" {
			next.ServeHTTP(w, r)
			return
		}

		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			stats.IncUnauthorized()
			writeJSON(w, http.StatusUnauthorized, map[string]any{
				"status": "unauthorized",
				"error":  "invalid or missing bearer token",
			})
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
		if token == "" {
			stats.IncUnauthorized()
			writeJSON(w, http.StatusUnauthorized, map[string]any{
				"status": "unauthorized",
				"error":  "invalid or missing bearer token",
			})
			return
		}

		if _, ok := ts.tokens[token]; !ok {
			stats.IncUnauthorized()
			writeJSON(w, http.StatusUnauthorized, map[string]any{
				"status": "unauthorized",
				"error":  "invalid or missing bearer token",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
