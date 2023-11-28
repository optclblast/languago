package api

import (
	"encoding/json"
	"languago/pkg/errors"
)

func (a *API) responseError(msg string, e error, code int) []byte {
	err := a.errorsPresenter.ServiceError(
		e,
		errors.ErrorServiceID(a.ID),
		errors.ErrorServiceErr(errors.New(errors.Code(code), msg)),
	)

	body, err := json.Marshal(a.errorsPresenter.ResponseError(err))
	if err != nil {
		a.log.Error("error responding to request: ", "error: ", err)
	}

	return body
}
