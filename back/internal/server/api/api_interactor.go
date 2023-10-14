package api

import (
	"encoding/json"
	"fmt"
	"languago/internal/pkg/logger"
	"languago/internal/pkg/models/requests/rest"
	"net/http"
)

func (a *FlashcardAPI) responseError(w http.ResponseWriter, e error, code int) error {
	a.log.Warn(logger.ErrorField, logger.LogFieldPair(logger.ErrorField, e))
	resp := rest.NewFlashcardResponse{
		Errors: []string{e.Error()},
	}
	body, err := json.Marshal(resp)
	if err != nil {
		a.log.Warn("error responding to request: ", logger.LogFieldPair(logger.ErrorField, err))
		return fmt.Errorf("error responding ro request: %w", err)
	}

	if e != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	_, err = w.Write(body)
	if err != nil {
		a.log.Warn("error responding to request: ", logger.LogFieldPair(logger.ErrorField, err))
		return fmt.Errorf("error responding ro request: %w", err)
	}
	return nil
}
