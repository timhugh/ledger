package middleware

import (
	"context"
	"github.com/google/uuid"
	"github.com/timhugh/ctxlogger"
	"net/http"
)

func Log(ctx context.Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.NewString()
		ctx := ctxlogger.AddParam(ctx, "request_id", requestID)
		r = r.WithContext(ctx)

		requestContext := ctxlogger.AddParam(ctxlogger.AddParam(ctx, "method", r.Method), "path", r.URL.Path) // gross
		ctxlogger.Info(requestContext, "")

		next.ServeHTTP(w, r)
	})
}
