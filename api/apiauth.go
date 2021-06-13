package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/sreesanthv/go-api-base/services"
)

func Authenticator(s *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())
			if err != nil {
				s.Logger.Error(err)
				sendAuthError(w, s, "Invalid access token")
				return
			}

			if token == nil || jwt.Validate(token) != nil {
				s.Logger.Error(err)
				sendAuthError(w, s, "Invalid access token")
				return
			}

			// validate claim against redis info
			id, _ := token.Get("user_id")
			userIdClaim, err := strconv.ParseInt(fmt.Sprintf("%v", id), 10, 32)
			id, _ = token.Get("access_uuid")
			idRedis, err := s.Redis.Get(id.(string))
			if err != nil {
				sendAuthError(w, s, "Access token has been expired")
				return
			}
			userIdRedis, err := strconv.ParseInt(fmt.Sprintf("%s", idRedis), 10, 32)
			if userIdClaim == 0 || userIdClaim != userIdRedis {
				sendAuthError(w, s, "Access token has been modified")
				return
			}

			//valid token
			next.ServeHTTP(w, r)
		})
	}
}

//send error response
func sendAuthError(w http.ResponseWriter, s *services.AuthService, message string) {
	dt := &responseData{
		Status:  "nok",
		Message: message,
	}

	jData, err := json.Marshal(dt)
	if err != nil {
		s.Logger.Error(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write(jData)
}
