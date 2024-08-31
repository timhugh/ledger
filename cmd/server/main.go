package main

import (
	"github.com/timhugh/ledger/cmd/server/app/journals"
	"github.com/timhugh/ledger/cmd/server/middleware"
	"github.com/timhugh/ledger/db/sqlite"
	"log"
	"net/http"
)

func main() {
	repo, err := sqlite.Open("development.db")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("GET /ping", wrap(Ping()))

	http.Handle("GET /api/journals/{journal_id}", wrapAPI(journals.GetJournalJson(repo)))
	http.Handle("GET /journals/{journal_id}", wrap(journals.GetJournalHtml(repo)))

	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func wrap(handlerFunc http.Handler) http.Handler {
	return middleware.RequestID(middleware.Log(handlerFunc))
}

func wrapAPI(handlerFunc http.Handler) http.Handler {
	return middleware.RequestID(middleware.Log(middleware.JSON(handlerFunc)))
}
