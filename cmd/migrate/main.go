package main

import (
	"context"
	"github.com/timhugh/ctxlogger"
	"github.com/timhugh/ledger/db/sqlite"
	"os"
)

func main() {
	ctx := context.Background()
	if len(os.Args) < 2 {
		ctxlogger.Error(ctx, "missing database file argument")
		os.Exit(1)
	}

	file := os.Args[1]
	ctxlogger.Info(ctx, "opening database %s", file)
	client, err := sqlite.Open(file)
	if err != nil {
		ctxlogger.Error(ctx, "failed to open database: %s", err.Error())
		os.Exit(1)
	}

	ctxlogger.Info(ctx, "running migrations\n")

	err = client.Migrate(ctx)
	if err != nil {
		ctxlogger.Error(ctx, "migration failed: %s", err.Error())
		os.Exit(1)
	}

	ctxlogger.Info(ctx, "migrations complete\n")
}
