package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Sumitwarrior7/social/internal/store"
	"github.com/golang-jwt/jwt/v5"
)

//   - After a successful login, the browser caches the Authorization header containing the credentials.
//   - For every subsequent request to the same domain, path, or resource that requires authentication, the browser automatically
//     attaches the cached Authorization header, so you don't need to re-enter your credentials.
func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read the auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedBasicError(w, r, fmt.Errorf("authorization header is missing"))
				return
			}

			// parse it -> get the base64
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedBasicError(w, r, fmt.Errorf("authorization header is malformed"))
				return
			}

			// decode it
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedBasicError(w, r, err)
				return
			}

			// check the credentials
			username := app.config.auth.basic.user
			pass := app.config.auth.basic.pass

			creds := strings.SplitN(string(decoded), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != pass {
				app.unauthorizedBasicError(w, r, fmt.Errorf("invalid credentials"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (app *application) TokenAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read the auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedError(w, r, fmt.Errorf("authorization header is missing"))
				return
			}

			// parse it -> get the base64
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				app.unauthorizedError(w, r, fmt.Errorf("authorization header is malformed"))
				return
			}

			// decode it
			token := parts[1]
			jwtToken, err := app.authenticator.ValidateToken(token)
			if err != nil {
				app.unauthorizedError(w, r, err)
				return
			}

			claims, _ := jwtToken.Claims.(jwt.MapClaims)
			userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
			if err != nil {
				app.unauthorizedError(w, r, err)
				return
			}

			ctx := r.Context()
			user, err := app.GetUser(ctx, userId)
			if err != nil {
				app.unauthorizedError(w, r, err)
				return
			}

			ctx = context.WithValue(ctx, userCtx, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Authorization
func (app *application) CheckPostOwnership(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUserFromCtx(r)
		post := getPostFromCtx(r)

		// If the post actually belongs to the user
		log.Println("User id :", user.Id)
		log.Println("Post User id :", post.UserId)
		if post.UserId == user.Id {
			next.ServeHTTP(w, r)
			return
		}

		// Role based precedence
		ctx := r.Context()
		allowed, err := app.checkRolePrecedence(ctx, user, requiredRole)
		if err != nil {
			app.internalServerError(w, r, err)
		}

		if !allowed {
			app.forbidenWarning(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) CheckCommentOwnership(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUserFromCtx(r)
		comment := getCommentFromCtx(r)
		// If the post actually belongs to the user
		if comment.UserId == user.Id {
			next.ServeHTTP(w, r)
			return
		}

		// Role based precedence
		ctx := r.Context()
		allowed, err := app.checkRolePrecedence(ctx, user, requiredRole)
		if err != nil {
			app.internalServerError(w, r, err)
		}

		if !allowed {
			app.forbidenWarning(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) checkRolePrecedence(ctx context.Context, user *store.User, roleName string) (bool, error) {
	role, err := app.store.Roles.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}

	return user.Role.Level >= role.Level, nil
}

func (app *application) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.rateLimiter.Enabled {
			if allow, retryAfter := app.rateLimiter.Allow(r.RemoteAddr); !allow {
				app.rateLimitExceededError(w, r, retryAfter.String())
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

/* Helper Function */
// Caching used
func (app *application) GetUser(ctx context.Context, userId int64) (*store.User, error) {
	// If redis is not enabled, then we will fetch directly from database
	if !app.config.redisCfg.enabled {
		return app.store.Users.GetById(ctx, userId)
	}

	// Checking wether user details is present in cache or not
	user, err := app.cacheStorage.Users.Get(ctx, userId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		log.Println("DB Hit!!!!!")
		// Retrieving user details from database
		user, err = app.store.Users.GetById(ctx, userId)
		if err != nil {
			return nil, err
		}

		// Storing user details in cache
		err := app.cacheStorage.Users.Set(ctx, user)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}
