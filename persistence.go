package ledger

type JournalRepository interface {
	CreateJournal(journal *Journal) error
	GetJournal(uuid string) (*Journal, error)
}
