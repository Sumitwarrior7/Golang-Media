package main

import (
	"log"
	"net/http"

	"github.com/Sumitwarrior7/social/internal/store"
)

// getUserFeedHandler godoc
//
//	@Summary		Fetches the user feed
//	@Description	Fetches the user feed
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Param			since	query		string	false	"Since"
//	@Param			until	query		string	false	"Until"
//	@Param			limit	query		int		false	"Limit"
//	@Param			offset	query		int		false	"Offset"
//	@Param			sort	query		string	false	"Sort"
//	@Param			tags	query		string	false	"Tags"
//	@Param			search	query		string	false	"Search"
//	@Success		200		{object}	[]store.PostWithMetadata
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	// Pagination, filter, sort
	fq := store.PaginatedFeedQuery{
		// Default Paginated Values
		Limit:  20,
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

	ctx := r.Context()
	user := getUserFromCtx(r)

	feeds, err := app.store.Posts.GetUserFeed(ctx, user.Id, fq)
	log.Println(feeds)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Iterate through feeds and fetch comments for each feed
	for i := range feeds {
		comments, err := app.store.Comments.GetByPostId(ctx, feeds[i].Id)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		// Add the fetched comments to the current feed
		feeds[i].Comments = comments
	}

	log.Println(feeds)

	if err := app.jsonResponse(w, http.StatusOK, feeds); err != nil {
		app.internalServerError(w, r, err)
	}
}
