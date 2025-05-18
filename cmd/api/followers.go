package main

import "net/http"

type followerPayload struct {
	FollowedId int64 `json:"followed_id"`
}

func (app *application) followHandler(w http.ResponseWriter, r *http.Request) {
	var payload = followerPayload{}

	if err := readJson(w, r.Body, &payload); err != nil {
		writeBadInputErr(w, err)
		return
	}

	follower := app.getUserFromContext(r)
	println(payload.FollowedId)
	println(follower.ID)

	err := app.store.Followers.Follow(r.Context(), payload.FollowedId, follower.ID)
	if err != nil {
		writeInternalServerErr(w, err)
		return
	}
	writeJson(w, http.StatusCreated, "ok")
}

func (app *application) unfollowHandler(w http.ResponseWriter, r *http.Request) {
	var payload = followerPayload{}

	if err := readJson(w, r.Body, &payload); err != nil {
		writeBadInputErr(w, err)
		return
	}

	follower := app.getUserFromContext(r)

	err := app.store.Followers.Unfollow(r.Context(), payload.FollowedId, follower.ID)
	if err != nil {
		writeInternalServerErr(w, err)
		return
	}
	writeJson(w, http.StatusCreated, "ok")
}
