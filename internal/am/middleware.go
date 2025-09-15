package am

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/csrf"
)

const defaultCSRFKey = "set-a-csrf-key!"

var (
	csrfMiddleware func(http.Handler) http.Handler
	once           sync.Once
)

// LogHeadersMw is a middleware that logs all request headers.
func LogHeadersMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := NewLogger("request-headers")
		log.Info("Incoming Request Headers:")
		for name, headers := range r.Header {
			for _, h := range headers {
				log.Infof("  %s: %s", name, h)
			}
		}
		next.ServeHTTP(w, r)
	})
}

// MethodOverrideMw is a middleware that checks for a _method form field and overrides the request method.
func MethodOverrideMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if override := r.FormValue("_method"); override != "" {
				r.Method = override
			}
		}
		next.ServeHTTP(w, r)
	})
}

// CSRFMw is a middleware that protects against CSRF attacks.
func CSRFMw(cfg *Config) func(next http.Handler) http.Handler {
	if cfg == nil {
		return passThroughMw
	}

	initCSRF(cfg)

	return func(next http.Handler) http.Handler {
		return csrfMiddleware(next)
	}
}

func passThroughMw(next http.Handler) http.Handler {
	return next
}

func initCSRF(cfg *Config) {
	once.Do(func() {
		key := cfg.StrValOrDef(Key.SecCSRFKey, defaultCSRFKey)
		to := cfg.StrValOrDef(Key.SecCSRFRedirect, "/csrf-error")

		csrfMiddleware = csrf.Protect(
			[]byte(key),
			csrf.FieldName(CSRFFieldName),
			csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, to, http.StatusFound)
			})),
		)
	})
}

// ReqIDKey is the context key for the request ID.
var ReqIDKey = contextKey("requestID")

// RequestIDMw is a middleware that assigns a unique ID to each request and stores it in the context and as a header.
func RequestIDMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.NewString()
		ctx := context.WithValue(r.Context(), ReqIDKey, id)
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ReqID returns the request ID from the context, or an empty string if not set.
func ReqID(r *http.Request) string {
	if v := r.Context().Value(ReqIDKey); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}

// UserKey is the context key for the authenticated user.
var UserKey = contextKey("user")

// SessionStore is a placeholder interface for session management.
type SessionStore interface {
	GetUserFromSession(ctx context.Context, sessionID string) (uuid.UUID, error)
}

// UserService is a placeholder interface for user retrieval.
type UserService interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*UserCtxData, error)
}

// EncryptionKeyMw is a middleware that injects the encryption key into the request context.
func EncryptionKeyMw(app *App) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := app.Cfg().ByteSliceVal(Key.SecEncryptionKey)
			if len(key) == 0 {
				app.Log().Error("Encryption key is empty in EncryptionKeyMw")
			} else {
				app.Log().Debug("Encryption key loaded in EncryptionKeyMw")
			}
			ctx := context.WithValue(r.Context(), EncryptionKeyCtxKey, key)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// InternalAuthMiddleware is a middleware that verifies the internal authentication token.
func InternalAuthMiddleware(app *App) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := app.Log()
			internalToken := r.Header.Get(InternalAuthHeader)
			if internalToken == "" || internalToken != app.InternalAuthToken() {
				log.Errorf("Unauthorized: Invalid or missing internal auth token.")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// AuthMw is a middleware that handles authentication for both internal and external requests.
func AuthMw(app *App, sessionStore SessionStore, userService UserService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := app.Log() // Get logger from app instance

			// Development bypass for authentication
			if app.Cfg().BoolVal(Key.SecBypassAuth, false) {
				log.Infof("Authentication bypassed for development purposes!")
				// Create a dummy user for development
				user := &UserCtxData{
					ID:                    uuid.NewSHA1(uuid.Nil, []byte("dev-user")),
					Permissions:           []uuid.UUID{uuid.NewSHA1(uuid.Nil, []byte("dev-perm"))}, // Placeholder for dev permission
					ContextualPermissions: make(map[uuid.UUID]uuid.UUID),
				}
				ctx := context.WithValue(r.Context(), UserKey, user)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// 1. Check for session cookie (for webapp requests)
			sessionCookie, err := r.Cookie("user_session") // Assuming "session_id" is the session cookie name
			if err == nil && sessionCookie != nil {
				userID, err := sessionStore.GetUserFromSession(r.Context(), sessionCookie.Value)
				if err == nil && userID != uuid.Nil {
					user, err := userService.GetUserByID(r.Context(), userID)
					if err == nil && user != nil {
						log.Infof("Webapp request authenticated via session cookie. User ID: %s", user.ID)
						ctx := context.WithValue(r.Context(), UserKey, user)
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}
				}
				log.Errorf("Invalid session cookie: %v", err)
			}

			// 2. If not session, check for JWT (for external requests)
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				jwtToken := strings.TrimPrefix(authHeader, "Bearer ")
				// TODO: Implement JWT validation here.
				// Validate jwtToken and extract user info.
				// If valid, authenticate user and set in context.
				log.Infof("External request with JWT: %s", jwtToken)
				user := &UserCtxData{
					ID:                    uuid.NewSHA1(uuid.Nil, []byte("jwt-user")),
					Permissions:           []uuid.UUID{uuid.NewSHA1(uuid.Nil, []byte("user-perm"))}, // Placeholder for user permission
					ContextualPermissions: make(map[uuid.UUID]uuid.UUID),
				}
				ctx := context.WithValue(r.Context(), UserKey, user)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// 3. If no valid authentication, return 401 Unauthorized
			log.Errorf("Unauthorized request: No valid session cookie or JWT found.")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		})
	}
}
