package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lutif/go-social/internal/store"
)

type updateCommentPayload struct {
	Content string `json:"content"`
	Likes   int64  `json:"likes"`
}

type createCommentPayload struct {
	Content string `json:"content"`
	PostID  int64  `json:"post_id"`
}
type commentKey string

const ctxCommentKey commentKey = "comment"

func (app *application) createComment(w http.ResponseWriter, r *http.Request) {
	post, _ := r.Context().Value(ctxPostKey).(store.Post)
	userID := 1
	payload := createCommentPayload{}
	err := readJson(w, r.Body, &payload)

	if err != nil {
		writeBadInputErr(w, err)
		return
	}

	comment := store.Comment{
		AuthorID: int64(userID),
		Content:  payload.Content,
		PostID:   post.ID,
	}

	err = app.store.Comments.Create(r.Context(), &comment)

	if err != nil {
		writeInternalServerErr(w, err)
		return
	}
	writeJson(w, http.StatusCreated, comment)
}

func (app *application) getComment(w http.ResponseWriter, r *http.Request) {
	comment := app.getCommentFromContext(r)
	writeJson(w, http.StatusAccepted, comment)
}

func (app *application) updateComment(w http.ResponseWriter, r *http.Request) {
	comment := app.getCommentFromContext(r)
	payload := updateCommentPayload{}

	if err := readJson(w, r.Body, &payload); err != nil {
		writeBadInputErr(w, err)
		return
	}

	comment.Content = payload.Content

	err := app.store.Comments.Update(r.Context(), &comment)
	if err != nil {
		writeNotFoundError(w, err)
		return
	}

	writeJson(w, http.StatusAccepted, comment)
}

func (app *application) deleteComment(w http.ResponseWriter, r *http.Request) {
	comment := app.getCommentFromContext(r)
	err := app.store.Comments.Delete(r.Context(), comment.ID)

	if err != nil {
		writeNotFoundError(w, err)
		return
	}
}

func (app *application) commentContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "commentId")
		commentId, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			writeBadInputErr(w, err)
			return
		}
		ctx := r.Context()
		comment, err := app.store.Comments.GetById(ctx, commentId)

		if err != nil {
			writeNotFoundError(w, err)
			return
		}
		ctx = context.WithValue(ctx, ctxCommentKey, comment)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) getCommentFromContext(r *http.Request) store.Comment {
	return r.Context().Value(ctxCommentKey).(store.Comment)
}
