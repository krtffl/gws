package middlewares

import (
	"context"
	"net/http"

	"github.com/krtffl/gws/internal/cookie"
	"github.com/krtffl/gws/internal/logger"
)

func Solved(
	cookieSvc *cookie.Service,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			logger.Info("[AuthMiddleware - Solved] " +
				"Checking if was solved")

			challenge, err := cookieSvc.Retrieve(r)
			if err != nil {
				logger.Info(
					"[AuthMiddleware - AuthenticateFrontend] "+
						"Couldn't retrieve session from request. %v",
					err,
				)
				http.Redirect(w, r, "/challenge", http.StatusFound)
				return
			}

			ctx := context.WithValue(r.Context(), "challenge", challenge)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
