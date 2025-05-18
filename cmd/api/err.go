package main

import (
	"errors"
	"net/http"

	"github.com/lutif/go-social/internal/store"
)

func writeInternalServerErr(w http.ResponseWriter, err error) {
	writeErrorJson(w, http.StatusInternalServerError, err.Error())
}

func writeBadInputErr(w http.ResponseWriter, err error) {
	writeErrorJson(w, http.StatusBadRequest, err.Error())
}

func writeNotFoundError(w http.ResponseWriter, err error) {
	if errors.Is(err, store.ErrNotFound) {
		writeErrorJson(w, http.StatusNotFound, err.Error())
		return
	}
	writeErrorJson(w, http.StatusInternalServerError, err.Error())
}
