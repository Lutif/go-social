package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lutif/go-social/internal/store"
)

type postKey string

const ctxPostKey postKey = "post"

type createPostPayload struct {
	Title   *string `json:"title" validate:"required,max=100"`
	Content *string `json:"content" validate:"required,max=1000"`
	Tags    *[]string
}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	var payload createPostPayload
	userId := int64(1)
	if err := readJson(w, r.Body, &payload); err != nil {
		writeBadInputErr(w, err)
		return
	}

	post := store.Post{
		Title:   *payload.Title,
		Content: *payload.Content,
		Tags:    *payload.Tags,
		UserID:  userId,
	}

	if err := app.store.Posts.Create(r.Context(), &post); err != nil {
		println(err.Error())
		writeNotFoundError(w, err)
		return
	}
	writeJson(w, http.StatusCreated, post)
}

func (app *application) getPost(w http.ResponseWriter, r *http.Request) {
	post, _ := r.Context().Value(ctxPostKey).(store.Post)
	writeJson(w, http.StatusOK, post)
}

func (app *application) getFeed(w http.ResponseWriter, r *http.Request) {
	userId := int64(1) // replace with auth
	var posts = []store.Post{}
	err := app.store.Posts.ListByUserId(r.Context(), userId, &posts)
	if err != nil {
		writeNotFoundError(w, err)
		return
	}
	writeJson(w, http.StatusOK, posts)

}

func (app *application) updatePost(w http.ResponseWriter, r *http.Request) {

	post, _ := r.Context().Value(ctxPostKey).(store.Post)
	payload := createPostPayload{}
	err := readJson(w, r.Body, &payload)
	if err != nil {
		writeInternalServerErr(w, err)
		return
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Title != nil {
		post.Title = *payload.Title
	}
	if payload.Tags != nil {
		post.Tags = *payload.Tags
	}

	err = app.store.Posts.Update(r.Context(), &post)
	if err != nil {
		writeNotFoundError(w, err)
		return
	}
	writeJson(w, http.StatusOK, post)
}

func (app *application) deletePost(w http.ResponseWriter, r *http.Request) {
	post, _ := r.Context().Value(ctxPostKey).(store.Post)

	if post.ID == 0 {
		writeNotFoundError(w, store.ErrNotFound)
		return
	}

	err := app.store.Posts.Delete(r.Context(), post.ID)
	if err != nil {
		writeNotFoundError(w, err)
	}
}

func (app *application) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postId")
		postId, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			writeBadInputErr(w, err)
			return
		}
		ctx := r.Context()
		post, err := app.store.Posts.GetById(ctx, postId)
		if err != nil {
			writeNotFoundError(w, err)
			return
		}
		println(post.Title)
		ctx = context.WithValue(ctx, ctxPostKey, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
