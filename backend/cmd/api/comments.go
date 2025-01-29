package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/Sumitwarrior7/social/internal/store"
	"github.com/go-chi/chi/v5"
)

type commentKey string

const commentCtx commentKey = "comment"

type CommentPayload struct {
	Content string `json:"content" validate:"required,max=1000"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload CommentPayload

	err := readJSON(w, r, &payload)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	user := getUserFromCtx(r)
	post := getPostFromCtx(r)
	comment := &store.Comment{
		PostId:  post.Id,
		UserId:  user.Id,
		Content: payload.Content,
	}

	ctx := r.Context()
	if err := app.store.Comments.Create(ctx, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getCommentByIdHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	comment := getCommentFromCtx(r)

	comment, err := app.store.Comments.GetById(ctx, comment.Id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	comment := getCommentFromCtx(r)

	if err := app.store.Comments.Delete(ctx, comment.Id); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (app *application) updateCommentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	comment := getCommentFromCtx(r)

	var payload CommentPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Updated Post Validations
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	updatedComment := &store.Comment{
		Id:      comment.Id,
		Content: payload.Content,
	}
	// if payload.Content != "" {
	// 	comment.Content = payload.Content
	// }

	if err := app.store.Comments.Update(ctx, updatedComment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, updatedComment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) commentContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "commentId")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()
		comment, err := app.store.Comments.GetById(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundError(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}
		ctx = context.WithValue(ctx, commentCtx, comment)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getCommentFromCtx(r *http.Request) *store.Comment {
	comment, _ := r.Context().Value(commentCtx).(*store.Comment)
	return comment
}
