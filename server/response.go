package server

import (
	"encoding/json"
	"net/http"
)

// This file contains functions to be used frequently for HTTP response
// delivery.

func sendResponse(w http.ResponseWriter, msg string, code int) error {

	encoded, err := json.Marshal(&genericResponse{Message: msg})
	if err != nil {
		return err
	}
	w.WriteHeader(code)
	_, err = w.Write(encoded)
	if err != nil {
		return err
	}
	return nil
}

func sendClientError(w http.ResponseWriter, msg string) error {
	return sendResponse(w, msg, http.StatusBadRequest)
}

func sendServerError(w http.ResponseWriter) error {
	return sendResponse(w, "server error", 500)
}

func sendCreated(w http.ResponseWriter) error {
	return sendResponse(w, "created", http.StatusCreated)
}

func sendUpdated(w http.ResponseWriter) error {
	return sendResponse(w, "updated", http.StatusOK)
}
