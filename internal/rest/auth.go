package rest

import (
	"net/http"
	"strings"
)

type TokenSet struct {
	enabled bool
	tokens  map[string]struct{} // token -> present
}

func NewTokenSet(enabled bool, tokenValues []string) *TokenSet {
	m := make(map[string]struct{}, len(tokenValues))
	for _, t := range tokenValues {
		if strings.TrimSpace(t) == "" {
			continue
		}
		m[t] = struct{}{}
	}
	return &TokenSet{
		enabled: enabled,
		tokens:  m,
	}
}

func (ts *TokenSet) Require(next http.Handler, stats *Stats) http.Handler {
	if !ts.enabled {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// health endpoint is always public
		if r.URL.Path == "/api/v1/health" {
			next.ServeHTTP(w, r)
			return
		}

		// ---- Authorization header check ----
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			stats.IncUnauthorized()
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
		if token == "" {
			stats.IncUnauthorized()
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if _, ok := ts.tokens[token]; !ok {
			stats.IncUnauthorized()
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// ---- authorized ----
		next.ServeHTTP(w, r)
	})
}
