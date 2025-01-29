package main

import (
	"net/http"
	"strconv"

	"github.com/Sumitwarrior7/social/internal/store"
	"github.com/go-chi/chi/v5"
)

type userKey string

const userCtx userKey = "user"

type FollowUser struct {
	UserId int64 `json:"user_id"`
}

// GetUser godoc
//
// @Summary     Fetches a user profile
// @Description Fetches a user profile by ID
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       id   path   int  true  "User ID"
// @Success     200  {object}  store.User
// @Failure     400  {object}  error
// @Failure     404  {object}  error
// @Failure     500  {object}  error
// @Security    ApiKeyAuth
// @Router      /users/{id} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "userId")
	userId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	user, err := app.GetUser(ctx, userId)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

// FollowUser godoc
//
// @Summary     Follows a user
// @Description Follows a user by ID
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       userID  path   int  true  "User ID"
// @Success     204     {string}  string  "User followed"
// @Failure     400     {object}  error   "User payload missing"
// @Failure     404     {object}  error   "User not found"
// @Security    ApiKeyAuth
// @Router      /users/{userID}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromCtx(r)
	followedId, err := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	if err := app.store.Followers.Follow(ctx, followerUser.Id, followedId); err != nil {
		switch err {
		case store.ErrConflict:
			app.conflictError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

// UnfollowUser gdoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User unfollowed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromCtx(r)
	unfollowedId, err := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	if err := app.store.Followers.Unfollow(ctx, followerUser.Id, unfollowedId); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	err := app.store.Users.Activate(r.Context(), token)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
}

// func (app *application) userContextMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		idParam := chi.URLParam(r, "userId")
// 		userId, err := strconv.ParseInt(idParam, 10, 64)
// 		if err != nil {
// 			app.badRequestError(w, r, err)
// 			return
// 		}
// 		// log.Println("User Id :%s", userId)

// 		ctx := r.Context()
// 		user, err := app.store.Users.GetById(ctx, userId)
// 		if err != nil {
// 			switch err {
// 			case store.ErrNotFound:
// 				app.notFoundError(w, r, err)
// 				return
// 			default:
// 				app.internalServerError(w, r, err)
// 				return
// 			}
// 		}

// 		ctx = context.WithValue(ctx, userCtx, user)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

func getUserFromCtx(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}

func (app *application) getCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getFollowedUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := getUserFromCtx(r)

	followedUsers, err := app.store.Followers.GetFollowedUsersById(ctx, user.Id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, followedUsers); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Pagination, filter, sort
	fq := store.PaginatedFeedQuery{
		Limit:  12,
		Offset: 0,
		Sort:   "desc",
	}
	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	//Validating PaginatedFeedQuery Struct
	if err := Validate.Struct(fq); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	allUsers, err := app.store.Users.GetAllUsers(ctx, fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, allUsers); err != nil {
		app.internalServerError(w, r, err)
	}
}
