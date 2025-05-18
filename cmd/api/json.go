package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(
		validator.WithRequiredStructEnabled(),
	)
}

func writeJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJson(w http.ResponseWriter, r io.ReadCloser, result any) error {
	max_bytes := 1_048_578
	r = http.MaxBytesReader(w, r, int64(max_bytes))

	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(result)

	if err != nil {
		return err
	}

	err = Validate.Struct(result)

	if err != nil {
		return err
	}

	return nil
}

func writeErrorJson(w http.ResponseWriter, status int, message string) {

	type envelop struct {
		Error  string
		Status int
	}

	writeJson(w, status, envelop{
		Error:  message,
		Status: status,
	})
}
