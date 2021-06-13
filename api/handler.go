package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/jwtauth"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/sreesanthv/go-api-base/database"
)

const SUCCESS_REPONSE int = 2

type Handler struct {
	logger    *logrus.Logger
	store     *database.Store
	redis     *database.Redis
	validator *validator.Validate
}

func NewHandler(logger *logrus.Logger, store *database.Store, redis *database.Redis) *Handler {
	return &Handler{
		logger:    logger,
		store:     store,
		redis:     redis,
		validator: validator.New(),
	}
}

// read request
func (h *Handler) parseJSONBody(r *http.Request, reqData interface{}) error {
	buf := new(bytes.Buffer)

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		h.logger.Error(err)
	}

	err = json.Unmarshal(buf.Bytes(), reqData)
	if err != nil {
		h.logger.Error(err)
	}

	return err
}

type responseData struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// send success repsone
func (h *Handler) sendResponse(w http.ResponseWriter, resData interface{}) {
	dt := &responseData{
		Status: "ok",
		Data:   resData,
	}

	jData, err := json.Marshal(dt)
	if err != nil {
		h.logger.Error(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) badDataResponse(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Invalid request data"
	}

	dt := &responseData{
		Status:  "nok",
		Message: message,
	}

	jData, err := json.Marshal(dt)
	if err != nil {
		h.logger.Error(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(jData)
}

func (h *Handler) ServerError(w http.ResponseWriter) {
	dt := &responseData{
		Status:  "nok",
		Message: "Failed to process the request",
	}

	jData, err := json.Marshal(dt)
	if err != nil {
		h.logger.Error(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(jData)
}

func (h *Handler) unAuthorized(w http.ResponseWriter, message string) {
	if message == "" {
		message = "You don't have the permission to perform this request"
	}

	dt := &responseData{
		Status:  "nok",
		Message: message,
	}

	jData, err := json.Marshal(dt)
	if err != nil {
		h.logger.Error(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write(jData)
}

func (h *Handler) getClaims(r *http.Request) (map[string]interface{}, error) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		h.logger.Error(err)
	}

	return claims, err
}

func (h *Handler) getUserId(r *http.Request) (int64, error) {
	claims, err := h.getClaims(r)
	if err != nil {
		return 0, err
	}

	userId, err := strconv.ParseInt(fmt.Sprintf("%v", claims["user_id"]), 10, 32)
	if err != nil {
		h.logger.Error(err)
		return 0, err
	}

	return userId, nil
}

func (h *Handler) getAuthUuid(r *http.Request) (string, error) {
	claims, err := h.getClaims(r)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", claims["access_uuid"]), nil
}
