package controllers

import (
	"encoding/json"
	"net/http"
)

type GenericApiError struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"message"`
}

func WriteError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	errorResponse := GenericApiError{
		ErrorCode:    "BAD_BODY",
		ErrorMessage: err.Error(),
	}

	json.NewEncoder(w).Encode(errorResponse)
}
