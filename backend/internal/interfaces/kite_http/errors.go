package interfaces

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func writeError(w http.ResponseWriter, status int, errCode, msg string, err error) {
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
	} else {
		log.Error().Msg(msg)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: errCode, Message: msg})
}
