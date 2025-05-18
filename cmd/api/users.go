package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lutif/go-social/internal/store"
)

type createUserPayload struct {
	Username string
	Password string
	Email    string
}

type userCtxType string

var userCtxKey userCtxType = "user"

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	var payload createUserPayload

	if err := readJson(w, r.Body, &payload); err != nil {
		writeErrorJson(w, http.StatusInternalServerError, err.Error())
		return
	}

	user := store.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
	}

	if err := app.store.Users.Create(r.Context(), &user); err != nil {
		writeErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJson(w, http.StatusCreated, user)
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := app.getUserFromContext(r)
	writeJson(w, http.StatusAccepted, user)
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		param := chi.URLParam(r, "userId")
		userID, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			writeBadInputErr(w, err)
			return
		}

		ctx := r.Context()
		user, err := app.store.Users.GetById(ctx, userID)

		if err != nil {
			writeNotFoundError(w, err)
			return
		}

		ctx = context.WithValue(ctx, userCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) getUserFromContext(r *http.Request) store.User {
	return r.Context().Value(userCtxKey).(store.User)
}
