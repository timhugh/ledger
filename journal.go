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
	UUID string

	Name         string
	Transactions []Transaction
}

type Transaction struct {
	UUID        string
	JournalUUID string

	Description string
	Memo        string
	LineItems   []LineItem
}

func (t *Transaction) Valid() error {
	var sum int
	for _, lineItem := range t.LineItems {
		sum = sum + lineItem.Amount
	}
	if sum != 0 {
		return ErrTransactionUnbalanced
	}
	return nil
}

type LineItem struct {
	UUID            string
	TransactionUUID string

	Date    string
	Amount  int
	Account string
	Status  LineItemStatus
}
