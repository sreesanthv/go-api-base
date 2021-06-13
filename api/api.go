package api

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/spf13/viper"
	"github.com/sreesanthv/go-api-base/database"
	"github.com/sreesanthv/go-api-base/logging"
)

func New() (*chi.Mux, error) {
	logger := logging.NewLogger()

	db, err := database.DBConn()
	if err != nil {
		logger.WithField("module", "database").Error(err)
		return nil, err
	}

	redis := database.NewRedis(logger)
	store := database.NewStore(db, logger)
	jwtAuth := jwtauth.New(jwt.SigningMethodHS256.Name, []byte(viper.GetString("jwt_secret_access")), nil)

	handler := NewHandler(logger, store, redis)
	auth := NewAuthHandler(handler)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.DefaultCompress)
	r.Use(middleware.Timeout(15 * time.Second))
	r.Use(jwtauth.Verifier(jwtAuth))

	r.Use(logging.NewStructuredLogger(logger))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	if viper.GetBool("allow_cors") {
		r.Use(corsConfig().Handler)
	}

	r.Mount("/auth", auth.Router())
	r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("status"))
	})

	return r, nil
}

func corsConfig() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           86400,
	})
}
