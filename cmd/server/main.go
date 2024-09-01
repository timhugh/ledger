package main

import (
	"context"
	"github.com/timhugh/ctxlogger"
	"github.com/timhugh/ledger/cmd/server/middleware"
	"github.com/timhugh/ledger/db/sqlite"
	"net/http"
	"os"
)

func main() {
	ctx := context.Background()

	db := "development.db"
	repo, err := sqlite.Open(db)
	if err != nil {
		ctxlogger.Error(ctx, "failed to open database '%s' %w", db, err)
		os.Exit(1)
	}

	http.Handle("GET /ping", wrap(ctx, Ping()))
	http.Handle("GET /{$}", wrap(ctx, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxlogger.Info(r.Context(), "redirecting to /transactions")
		http.Redirect(w, r, "/transactions", http.StatusFound)
	})))

	ctxlogger.Info(ctx, "server listening on localhost:8080")
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		ctxlogger.Error(ctx, "server failed with error: %s", err.Error())
	}
}

func wrap(ctx context.Context, handler http.Handler) http.Handler {
	return middleware.Log(
		ctx,
		handler,
	)
}

func wrapAPI(ctx context.Context, handler http.Handler) http.Handler {
	return wrap(ctx, middleware.JSON(handler))
}
