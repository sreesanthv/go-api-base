package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sreesanthv/go-api-base/database"
	"github.com/sreesanthv/go-api-base/services"
)

type AuthHandler struct {
	Handler
	authService *services.AuthService
}

func NewAuthHandler(handler *Handler) *AuthHandler {
	ah := &AuthHandler{
		Handler:     *handler,
		authService: services.NewAuthService(handler.logger, handler.store, handler.redis),
	}

	return ah
}

func (h *AuthHandler) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/login", h.login)
	r.Post("/refresh", h.refreshToken)
	return r
}

type loginPostBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// login with email and password
func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {

	body := &loginPostBody{}
	err := h.parseJSONBody(r, body)

	if err != nil {
		h.badDataResponse(w, "")
		return
	}

	user := h.authService.GetUser(body.Email)
	if user.ID == 0 || h.authService.IsValidPassword(user, body.Password) == false {
		h.badDataResponse(w, "Invalid email & password combination")
		return
	}

	tokens, err := h.newToken(user)
	if err != nil {
		h.ServerError(w)
		return
	}

	h.sendResponse(w, tokens)
}

// create new token for user
func (h *AuthHandler) newToken(user *database.AccountStore) (map[string]string, error) {
	token, err := h.authService.CreateToken(user)
	if err != nil {
		return nil, err
	}

	err = h.authService.PersistToken(user.ID, token)
	if err != nil {
		return nil, err
	}

	tokens := map[string]string{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	}

	return tokens, nil
}

// create new tokens using refresh token
func (h *AuthHandler) refreshToken(w http.ResponseWriter, r *http.Request) {
	body := map[string]string{}
	err := h.parseJSONBody(r, &body)
	if err != nil {
		h.badDataResponse(w, "")
		return
	}

	reToken, ok := body["refresh_token"]
	if !ok || reToken == "" {
		h.badDataResponse(w, "")
		return
	}

	tInfo, err := h.authService.ParseToken(reToken, services.TokenTypeRefresh)
	if err != nil {
		h.unAuthorized(w, "Invalid refresh token")
		return
	}

	// drop existing refresh token
	err = h.authService.DropToken(tInfo.Uuid)
	if err != nil {
		h.ServerError(w)
		return
	}

	user := h.authService.GetUserById(tInfo.UserId)
	if user.ID == 0 {
		h.unAuthorized(w, "User doesn't exists")
		return
	}

	tokens, err := h.newToken(user)
	if err != nil {
		h.ServerError(w)
		return
	}

	h.sendResponse(w, tokens)
}
