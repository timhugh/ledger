package journals

import (
	"context"
	"github.com/timhugh/ledger"
	"log"
	"net/http"
)

func logError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func logWriteError(_ int, err error) {
	if err != nil {
		log.Println(err)
	}
}

func GetJournalHtml(repo ledger.JournalGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		journalID := r.PathValue("journal_id")
		journal, err := repo.GetJournal(journalID)
		if err != nil {
			logWriteError(w.Write([]byte(err.Error())))
			return
		}
		view := GetJournalView(journal)
		logError(view.Render(context.TODO(), w))
	}
}

func GetJournalJson(repo ledger.JournalGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		journalID := r.PathValue("journal_id")
		journal, err := repo.GetJournal(journalID)
		if err != nil {
			logWriteError(w.Write([]byte(err.Error())))
			return
		}
		log.Println(journal)
	}
}
