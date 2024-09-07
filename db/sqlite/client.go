package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/timhugh/ledger"
	"strings"
)

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
var sessionColumns = []string{
	"sessions.session_uuid as session_uuid",
	"sessions.user_uuid as session_user_uuid",
}
var userColumns = []string{
	"users.user_uuid as user_uuid",
	"users.login as user_login",
	"users.password_hash as user_password_hash",
	"users.password_salt as user_password_salt",
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

func (c *Client) CreateJournal(ctx context.Context, journal *ledger.Journal) error {
	if journal.JournalUUID == "" {
		journal.JournalUUID = uuid.NewString()
	}
	query := "INSERT INTO journals (journal_uuid, name) VALUES (?, ?)"
	_, err := c.db.ExecContext(ctx, query, journal.JournalUUID, journal.Name)
	return err
}

func (c *Client) GetJournal(ctx context.Context, uuid string) (*ledger.Journal, error) {
	query := fmt.Sprintf("SELECT %s FROM journals WHERE journals.journal_uuid = ?",
		strings.Join(journalColumns, ", "))

	var journal ledger.Journal
	if err := c.db.GetContext(ctx, &journal, query, uuid); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ledger.ErrNotFound
		}
		return nil, err
	}

	transactions, err := c.GetTransactions(ctx, journal.JournalUUID)
	if err != nil {
		return nil, err
	}
	journal.Transactions = transactions

	return &journal, nil
}

func (c *Client) GetTransactions(ctx context.Context, journalUUID string) ([]*ledger.Transaction, error) {
	query := fmt.Sprintf("SELECT %s FROM transactions JOIN transaction_line_items ON transactions.transaction_uuid = transaction_line_items.transaction_uuid WHERE transactions.journal_uuid = ?",
		strings.Join(transactionColumns, ", ")+", "+strings.Join(transactionLineItemColumns, ", "))

	rows, err := c.db.QueryxContext(ctx, query, journalUUID)
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

func (c *Client) GetSession(ctx context.Context, uuid string) (*ledger.Session, error) {
	query := fmt.Sprintf("SELECT %s, %s FROM sessions JOIN users ON sessions.user_uuid = users.user_uuid WHERE sessions.session_uuid = ?",
		strings.Join(sessionColumns, ", "),
		strings.Join(userColumns, ", "))
	var result struct {
		ledger.Session
		ledger.User
	}
	if err := c.db.GetContext(ctx, &result, query, uuid); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ledger.ErrNotFound
		}
		return nil, err
	}
	session := &result.Session
	session.User = &result.User
	return session, nil
}
