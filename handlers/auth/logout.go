package auth

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
)

type LogoutResponse struct {
	Message string `json:"message"`
}

var revokedTokens = struct {
	sync.Mutex
	tokens map[string]bool
}{tokens: make(map[string]bool)}

func RevokeToken(token string) {
	revokedTokens.Lock()
	defer revokedTokens.Unlock()
	revokedTokens.tokens[token] = true
}

func IsTokenRevoked(token string) bool {
	revokedTokens.Lock()
	defer revokedTokens.Unlock()
	return revokedTokens.tokens[token]
}

func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusBadRequest)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		RevokeToken(tokenString)

		response := LogoutResponse{
			Message: "Logout successful",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
