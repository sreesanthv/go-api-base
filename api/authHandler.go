package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-pg/pg"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	Handler
}

func NewAuthHandler(db *pg.DB, logger *logrus.Logger) *AuthHandler {
	ah := &AuthHandler{}
	ah.DB = db
	ah.Logger = logger
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
		h.badDataResponse(w)
		return
	}

	// TODO login logic
}
