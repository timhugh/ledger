package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/timhugh/ledger"
	"strings"
)

var (
	journalColumns     []string = []string{"l.uuid", "l.name"}
	transactionColumns []string = []string{"t.uuid", "t.journal_uuid", "t.description", "t.memo"}
	lineItemColumns    []string = []string{"li.uuid", "li.transaction_uuid", "li.date", "li.amount", "li.account", "li.status"}
)

type Client struct {
	db *sql.DB
}

func Open(path string) (*Client, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return &Client{db: db}, nil
}

func (c *Client) Close() error {
	return c.db.Close()
}

func (c *Client) CreateJournal(journal *ledger.Journal) error {
	if journal.UUID == "" {
		journal.UUID = uuid.NewString()
	}
	query := `INSERT INTO journals (uuid, name) VALUES (?, ?)`
	_, err := c.db.Exec(query, journal.UUID, journal.Name)
	return err
}

func (c *Client) GetJournal(uuid string) (*ledger.Journal, error) {
	query := fmt.Sprintf(`SELECT %s, %s, %s 
    FROM journals l
    LEFT JOIN transactions t ON l.uuid = t.journal_uuid
    LEFT JOIN transaction_line_items li ON t.uuid = li.transaction_uuid
    `,
		strings.Join(journalColumns, ", "),
		strings.Join(transactionColumns, ", "),
		strings.Join(lineItemColumns, ", "))

	result, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var currentJournal ledger.Journal
	var currentTransaction ledger.Transaction
	for result.Next() {
		var transactionUUID sql.NullString
		var transactionJournalUUID sql.NullString
		var transactionDescription sql.NullString
		var transactionMemo sql.NullString

		var lineItemUUID sql.NullString
		var lineItemTransactionUUID sql.NullString
		var lineItemDate sql.NullString
		var lineItemAmount sql.NullInt64
		var lineItemAccount sql.NullString
		var lineItemStatus sql.NullString

		err = result.Scan(
			&currentJournal.UUID,
			&currentJournal.Name,
			&transactionUUID,
			&transactionJournalUUID,
			&transactionDescription,
			&transactionMemo,
			&lineItemUUID,
			&lineItemTransactionUUID,
			&lineItemDate,
			&lineItemAmount,
			&lineItemAccount,
			&lineItemStatus,
		)
		if err != nil {
			return nil, err
		}

		if transactionUUID.Valid {
			if currentTransaction.UUID == "" {
				currentTransaction = ledger.Transaction{
					UUID:        transactionUUID.String,
					JournalUUID: currentJournal.UUID,
					Description: transactionDescription.String,
					Memo:        transactionMemo.String,
				}
			}
			if currentTransaction.UUID != transactionUUID.String {
				currentJournal.Transactions = append(currentJournal.Transactions, currentTransaction)
				currentTransaction = ledger.Transaction{
					UUID:        transactionUUID.String,
					JournalUUID: currentJournal.UUID,
					Description: transactionDescription.String,
					Memo:        transactionMemo.String,
				}
			}
		} else {
			continue
		}

		if lineItemUUID.Valid {
			currentTransaction.LineItems = append(currentTransaction.LineItems, ledger.LineItem{
				UUID:            lineItemUUID.String,
				TransactionUUID: lineItemTransactionUUID.String,
				Date:            lineItemDate.String,
				Amount:          int(lineItemAmount.Int64),
				Account:         lineItemAccount.String,
				Status:          ledger.LineItemStatus(lineItemStatus.String),
			})
		}
	}
	currentJournal.Transactions = append(currentJournal.Transactions, currentTransaction)
	return &currentJournal, nil
}
