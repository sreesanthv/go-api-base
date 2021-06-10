package api

import (
	"net/http"

	"github.com/go-chi/chi"
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
	r.Post("/", h.login)
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

	token, err := h.authService.CreateToken(user)
	if err != nil {
		h.ServerError(w)
		return
	}

	err = h.authService.PersistToken(user.ID, token)
	if err != nil {
		h.ServerError(w)
		return
	}

	tokens := map[string]string{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	}

	h.sendResponse(w, tokens)
}
