package interfaces

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func writeError(ctx context.Context, w http.ResponseWriter, status int, errCode, msg string, err error) {
	logger := log.Ctx(ctx)

	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
	} else {
		logger.Error().Msg(msg)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: errCode, Message: msg})
}
