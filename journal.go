package ledger

import (
	"errors"
)

var (
	ErrTransactionUnbalanced = errors.New("transaction is unbalanced")
)

type LineItemStatus string

const (
	Cleared LineItemStatus = "cleared"
	Pending LineItemStatus = "pending"
)

type Journal struct {
	JournalUUID string `json:"journal_uuid" db:"journal_uuid"`

	Name         string         `json:"name" db:"journal_name"`
	Transactions []*Transaction `json:"transactions" db:"transactions"`
}

type JournalCreator interface {
	CreateJournal(journal *Journal) error
}
type JournalGetter interface {
	GetJournal(uuid string) (*Journal, error)
}

type Transaction struct {
	TransactionUUID string `json:"transaction_uuid" db:"transaction_uuid"`
	JournalUUID     string `json:"journal_uuid" db:"transaction_journal_uuid"`

	Description          string                 `json:"description" db:"transaction_description"`
	Memo                 string                 `json:"memo" db:"transaction_memo"`
	TransactionLineItems []*TransactionLineItem `json:"transaction_line_items" db:"transaction_line_items"`
}

type TransactionLister interface {
	GetTransactions(journalUUID string) ([]*Transaction, error)
}

func (t *Transaction) Valid() error {
	var sum int
	for _, lineItem := range t.TransactionLineItems {
		sum = sum + lineItem.Amount
	}
	if sum != 0 {
		return ErrTransactionUnbalanced
	}
	return nil
}

type TransactionLineItem struct {
	TransactionLineItemUUID string `json:"transaction_line_item_uuid" db:"transaction_line_item_uuid"`
	TransactionUUID         string `json:"transaction_uuid" db:"transaction_line_item_transaction_uuid"`

	Date    string         `json:"date" db:"transaction_line_item_date"`
	Amount  int            `json:"amount" db:"transaction_line_item_amount"`
	Account string         `json:"account" db:"transaction_line_item_account"`
	Status  LineItemStatus `json:"status" db:"transaction_line_item_status"`
}
