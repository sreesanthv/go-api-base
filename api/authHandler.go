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
	r.Post("/register", h.register)
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

	user := h.authService.GetAccount(body.Email)
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

	user := h.authService.GetAccountById(tInfo.UserId)
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

type userInfo struct {
	Name     string `validate:"required"`
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

// user registration
func (h *AuthHandler) register(w http.ResponseWriter, r *http.Request) {
	user := new(userInfo)
	err := h.parseJSONBody(r, user)
	if err != nil {
		h.badDataResponse(w, "")
		return
	}

	err = h.validator.Struct(user)
	if err != nil {
		h.badDataResponse(w, err.Error())
		return
	}

	// check email already used
	extAct := h.authService.GetAccount(user.Email)
	if extAct.ID != 0 {
		h.badDataResponse(w, "Email was already taken. Please try resetting password.")
		return
	}

	account, err := h.authService.CreateAccount(map[string]string{
		"name":     user.Name,
		"email":    user.Email,
		"password": user.Password,
	})
	if err != nil {
		h.ServerError(w)
		return
	}

	tokens, err := h.newToken(account)
	if err != nil {
		h.ServerError(w)
		return
	}

	h.sendResponse(w, tokens)
}
