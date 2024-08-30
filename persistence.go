package ledger

type JournalCreator interface {
	CreateJournal(journal *Journal) error
}
type JournalGetter interface {
	GetJournal(uuid string) (*Journal, error)
}
type TransactionLister interface {
	GetTransactions(journalUUID string) ([]*Transaction, error)
}
