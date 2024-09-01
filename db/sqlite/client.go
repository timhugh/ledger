package sqlite

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/timhugh/ledger"
	"strings"
)

var NoRecordError = errors.New("no record found")

var journalColumns = []string{
	"journals.journal_uuid as journal_uuid",
	"journals.name as journal_name",
}
var transactionColumns = []string{
	"transactions.transaction_uuid as transaction_uuid",
	"transactions.journal_uuid as transaction_journal_uuid",
	"transactions.description as transaction_description",
	"transactions.memo as transaction_memo",
}
var transactionLineItemColumns = []string{
	"transaction_line_items.transaction_line_item_uuid as transaction_line_item_uuid",
	"transaction_line_items.transaction_uuid as transaction_line_item_transaction_uuid",
	"transaction_line_items.date as transaction_line_item_date",
	"transaction_line_items.amount as transaction_line_item_amount",
	"transaction_line_items.account as transaction_line_item_account",
	"transaction_line_items.status as transaction_line_item_status",
}

type Client struct {
	db *sqlx.DB
}

func Open(path string) (*Client, error) {
	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return &Client{db: db}, nil
}

func (c *Client) Close() error {
	return c.db.Close()
}

func (c *Client) CreateJournal(journal *ledger.Journal) error {
	if journal.JournalUUID == "" {
		journal.JournalUUID = uuid.NewString()
	}
	query := `INSERT INTO journals (journal_uuid, name) VALUES (?, ?)`
	_, err := c.db.Exec(query, journal.JournalUUID, journal.Name)
	return err
}

func (c *Client) GetJournal(uuid string) (*ledger.Journal, error) {
	query := fmt.Sprint("SELECT ",
		strings.Join(journalColumns, ", "), " ",
		"FROM journals ",
		"WHERE journals.journal_uuid = ?")

	var journals []*ledger.Journal
	if err := c.db.Select(&journals, query, uuid); err != nil {
		return nil, err
	}

	if len(journals) == 0 {
		return nil, NoRecordError
	}

	journal := journals[0]
	transactions, err := c.GetTransactions(journal.JournalUUID)
	if err != nil {
		return nil, err
	}
	journal.Transactions = transactions

	return journal, nil
}

func (c *Client) GetTransactions(journalUUID string) ([]*ledger.Transaction, error) {
	query := fmt.Sprint("SELECT ",
		strings.Join(transactionColumns, ", "), ", ",
		strings.Join(transactionLineItemColumns, ", "), " ",
		"FROM transactions ",
		"JOIN transaction_line_items ON transactions.transaction_uuid = transaction_line_items.transaction_uuid ",
		"WHERE transactions.journal_uuid = ?")

	rows, err := c.db.Queryx(query, journalUUID)
	if err != nil {
		return nil, err
	}

	// TODO: map makes transaction order non-deterministic
	transactionsByID := make(map[string]*ledger.Transaction)
	for rows.Next() {
		var row struct {
			ledger.Transaction
			ledger.TransactionLineItem
		}

		err := rows.StructScan(&row)
		if err != nil {
			return nil, err
		}

		if transactionsByID[row.Transaction.TransactionUUID] == nil {
			transactionsByID[row.Transaction.TransactionUUID] = &row.Transaction
		}
		if row.TransactionLineItem.TransactionLineItemUUID != "" {
			transactionsByID[row.TransactionLineItem.TransactionUUID].TransactionLineItems =
				append(transactionsByID[row.TransactionLineItem.TransactionUUID].TransactionLineItems, &row.TransactionLineItem)
		}
	}

	var transactions []*ledger.Transaction
	for _, transaction := range transactionsByID {
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (c *Client) GetSession(uuid string) (*ledger.Session, error) {
    query := `SELECT session_uuid, user_uuid FROM sessions WHERE session_uuid = ?`
    var session ledger.Session
    if err := c.db.Get(&session, query, uuid); err != nil {
        return nil, err
    }
    return &session, nil
}
