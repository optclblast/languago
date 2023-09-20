package api

import (
	"encoding/json"
	"fmt"
	"languago/internal/pkg/models/requests/rest"
	"net/http"
)

func (a *API) response(w http.ResponseWriter, e error) error {
	resp := rest.NewFlashcardResponse{
		Errors: []string{e.Error()},
	}
	body, err := json.Marshal(resp)
	if err != nil {
		a.log.Warn("error responding to request: ", err)
		return fmt.Errorf("error responding ro request: %w", err)
	}

	if e != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	_, err = w.Write(body)
	if err != nil {
		a.log.Warn("error responding to request: ", err)
		return fmt.Errorf("error responding ro request: %w", err)
	}
	return nil
}
