package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/go-pg/pg"
	"github.com/sirupsen/logrus"
)

const SUCCESS_REPONSE int = 2

type Handler struct {
	DB     *pg.DB
	Logger *logrus.Logger
}

// read request
func (h *Handler) parseJSONBody(r *http.Request, reqData interface{}) error {
	buf := new(bytes.Buffer)

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		h.Logger.Error(err)
	}

	err = json.Unmarshal(buf.Bytes(), reqData)
	if err != nil {
		h.Logger.Error(err)
	}

	return err
}

type responseData struct {
	Status  string      `json:"status"`
	Message string      `json:"message",omitempty`
	Data    interface{} `json:"data"`
}

// send success repsone
func (h *Handler) sendResponse(w http.ResponseWriter, resData interface{}) {
	dt := &responseData{
		Status: "ok",
		Data:   resData,
	}

	jData, err := json.Marshal(dt)
	if err != nil {
		h.Logger.Error(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) badDataResponse(w http.ResponseWriter) {
	dt := &responseData{
		Status:  "nok",
		Message: "Invalid request data",
	}

	jData, err := json.Marshal(dt)
	if err != nil {
		h.Logger.Error(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
	w.WriteHeader(http.StatusBadRequest)
}
