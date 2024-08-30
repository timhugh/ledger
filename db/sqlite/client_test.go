package sqlite

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/matryer/is"
	_ "github.com/mattn/go-sqlite3"
	"github.com/timhugh/ledger"
	"testing"
)

var expectedJournalUUID = uuid.NewString()

const expectedJournalName = "Test Journal"

var expectedTransactionUUID = uuid.NewString()
var expectedTransactionJournalUUID = expectedJournalUUID

const expectedTransactionDescription = "Test Transaction"
const expectedTransactionMemo = "Test Memo"

var expectedTransactionUUID2 = uuid.NewString()
var expectedTransactionJournalUUID2 = expectedJournalUUID

const expectedTransactionDescription2 = "Test Transaction 2"
const expectedTransactionMemo2 = "Test Memo 2"

var expectedLineItemUUID = uuid.NewString()
var expectedLineItemTransactionUUID = expectedTransactionUUID

const expectedLineItemDate = "2024-08-01"
const expectedLineItemAmount = 100
const expectedLineItemAccount = "Expenses:Test"
const expectedLineItemStatus = ledger.Cleared

var expectedLineItemUUID2 = uuid.NewString()
var expectedLineItemTransactionUUID2 = expectedTransactionUUID

const expectedLineItemDate2 = "2024-08-02"
const expectedLineItemAmount2 = 200
const expectedLineItemAccount2 = "Expenses:Test2"
const expectedLineItemStatus2 = ledger.Cleared

var expectedLineItemUUID3 = uuid.NewString()
var expectedLineItemTransactionUUID3 = expectedTransactionUUID2

const expectedLineItemDate3 = "2024-08-03"
const expectedLineItemAmount3 = 300
const expectedLineItemAccount3 = "Expenses:Test3"
const expectedLineItemStatus3 = ledger.Pending

func testClient(t *testing.T) (*Client, *sqlx.DB) {
	is := is.New(t)

	db, err := sqlx.Open("sqlite3", ":memory:")
	is.NoErr(err)

	client := &Client{db: db}
	is.NoErr(client.Migrate())

	return client, db
}

func TestCreateJournal(t *testing.T) {
	is := is.New(t)
	t.Run("satisfies JournalCreator interface", func(t *testing.T) {
		client, _ := testClient(t)
		_, ok := interface{}(client).(ledger.JournalCreator)
		is.True(ok)
	})
	t.Run("creates journals", func(t *testing.T) {
		client, db := testClient(t)

		newJournal := &ledger.Journal{JournalUUID: expectedJournalUUID, Name: expectedJournalName}
		err := client.CreateJournal(newJournal)
		is.NoErr(err)

		var actualUUID string
		var actualName string
		err = db.QueryRow("SELECT journal_uuid, name FROM journals").Scan(&actualUUID, &actualName)
		is.NoErr(err)

		is.Equal(expectedJournalUUID, actualUUID)
		is.Equal(expectedJournalName, actualName)
	})
	t.Run("automatically assigns IDs", func(t *testing.T) {
		client, db := testClient(t)

		newJournal := &ledger.Journal{Name: expectedJournalName}
		is.NoErr(client.CreateJournal(newJournal))

		var journalUUID string
		var journalName string
		is.NoErr(db.QueryRow("SELECT journal_uuid, name FROM journals").Scan(&journalUUID, &journalName))

		if journalUUID == "" {
			t.Error("got empty uuid, want non-empty")
		}
		is.Equal(expectedJournalName, journalName)
	})
	t.Run("does not allow duplicate uuids", func(t *testing.T) {
		client, _ := testClient(t)

		existingJournal := &ledger.Journal{JournalUUID: expectedJournalUUID, Name: expectedJournalName}
		is.NoErr(client.CreateJournal(existingJournal))

		newJournal := &ledger.Journal{JournalUUID: expectedJournalUUID, Name: expectedJournalName}
		err := client.CreateJournal(newJournal)
		if err == nil {
			t.Error("expected error")
		}
	})
	t.Run("requires a name", func(t *testing.T) {
		// skipping for the moment -- go is passing the zero value of a string "" instead of nil
		// so the not null constraint is being satisfied
		t.Skip("not implemented")
		client, _ := testClient(t)

		actualJournal := &ledger.Journal{}
		err := client.CreateJournal(actualJournal)
		if err == nil {
			t.Error("expected error")
		}
	})
}
func TestGetJournal(t *testing.T) {
	is := is.New(t)
	t.Run("satisfies JournalGetter interface", func(t *testing.T) {
		client, _ := testClient(t)
		_, ok := interface{}(client).(ledger.JournalGetter)
		is.True(ok)
	})
	t.Run("returns a journal by id", func(t *testing.T) {
		client, db := testClient(t)

		_, err := db.Exec("INSERT INTO journals (journal_uuid, name) VALUES (?, ?)", expectedJournalUUID, expectedJournalName)
		is.NoErr(err)

		actualJournal, err := client.GetJournal(expectedJournalUUID)
		is.NoErr(err)
		is.Equal(expectedJournalUUID, actualJournal.JournalUUID)
		is.Equal(expectedJournalName, actualJournal.Name)
	})
	t.Run("returns a journal with transactions and line items", func(t *testing.T) {
		client, db := testClient(t)

		_, err := db.Exec("INSERT INTO journals (journal_uuid, name) VALUES (?, ?)", expectedJournalUUID, expectedJournalName)
		is.NoErr(err)
		_, err = db.Exec("INSERT INTO transactions (transaction_uuid, journal_uuid, description, memo) VALUES (?, ?, ?, ?)", expectedTransactionUUID, expectedJournalUUID, expectedTransactionDescription, expectedTransactionMemo)
		is.NoErr(err)
		_, err = db.Exec("INSERT INTO transaction_line_items (transaction_line_item_uuid, transaction_uuid, date, amount, account, status) VALUES (?, ?, ?, ?, ?, ?)", expectedLineItemUUID, expectedTransactionUUID, expectedLineItemDate, expectedLineItemAmount, expectedLineItemAccount, expectedLineItemStatus)
		is.NoErr(err)

		actualJournal, err := client.GetJournal(expectedJournalUUID)
		is.NoErr(err)
		is.Equal(expectedJournalUUID, actualJournal.JournalUUID)
		is.Equal(expectedJournalName, actualJournal.Name)
		is.Equal(1, len(actualJournal.Transactions))

		transaction := actualJournal.Transactions[0]
		is.Equal(expectedTransactionUUID, transaction.TransactionUUID)
		is.Equal(expectedTransactionJournalUUID, transaction.JournalUUID)
		is.Equal(expectedTransactionDescription, transaction.Description)
		is.Equal(expectedTransactionMemo, transaction.Memo)
		is.Equal(1, len(transaction.TransactionLineItems))

		transactionLineItem := transaction.TransactionLineItems[0]
		is.Equal(expectedLineItemUUID, transactionLineItem.TransactionLineItemUUID)
		is.Equal(expectedLineItemTransactionUUID, transactionLineItem.TransactionUUID)
		is.Equal(expectedLineItemDate, transactionLineItem.Date)
		is.Equal(expectedLineItemAmount, transactionLineItem.Amount)
		is.Equal(expectedLineItemAccount, transactionLineItem.Account)
		is.Equal(expectedLineItemStatus, transactionLineItem.Status)
	})
}
func TestGetTransactions(t *testing.T) {
	is := is.New(t)
	t.Run("satisfies TransactionLister interface", func(t *testing.T) {
		client, _ := testClient(t)
		_, ok := interface{}(client).(ledger.TransactionLister)
		is.True(ok)
	})
	t.Run("returns transactions for a journal", func(t *testing.T) {
		client, db := testClient(t)

		_, err := db.Exec("INSERT INTO journals (journal_uuid, name) VALUES (?, ?)", expectedJournalUUID, expectedJournalName)
		is.NoErr(err)
		_, err = db.Exec("INSERT INTO transactions (transaction_uuid, journal_uuid, description, memo) VALUES (?, ?, ?, ?)", expectedTransactionUUID, expectedJournalUUID, expectedTransactionDescription, expectedTransactionMemo)
		is.NoErr(err)
		_, err = db.Exec("INSERT INTO transaction_line_items (transaction_line_item_uuid, transaction_uuid, date, amount, account, status) VALUES (?, ?, ?, ?, ?, ?)", expectedLineItemUUID, expectedTransactionUUID, expectedLineItemDate, expectedLineItemAmount, expectedLineItemAccount, expectedLineItemStatus)
		is.NoErr(err)
		_, err = db.Exec("INSERT INTO transaction_line_items (transaction_line_item_uuid, transaction_uuid, date, amount, account, status) VALUES (?, ?, ?, ?, ?, ?)", expectedLineItemUUID2, expectedTransactionUUID, expectedLineItemDate2, expectedLineItemAmount2, expectedLineItemAccount2, expectedLineItemStatus2)
		is.NoErr(err)
		_, err = db.Exec("INSERT INTO transactions (transaction_uuid, journal_uuid, description, memo) VALUES (?, ?, ?, ?)", expectedTransactionUUID2, expectedJournalUUID, expectedTransactionDescription2, expectedTransactionMemo2)
		is.NoErr(err)
		_, err = db.Exec("INSERT INTO transaction_line_items (transaction_line_item_uuid, transaction_uuid, date, amount, account, status) VALUES (?, ?, ?, ?, ?, ?)", expectedLineItemUUID3, expectedTransactionUUID2, expectedLineItemDate3, expectedLineItemAmount3, expectedLineItemAccount3, expectedLineItemStatus3)
		is.NoErr(err)

		actualTransactions, err := client.GetTransactions(expectedJournalUUID)
		is.NoErr(err)

		is.Equal(2, len(actualTransactions))
		transaction1 := actualTransactions[0]
		is.Equal(expectedTransactionUUID, transaction1.TransactionUUID)
		is.Equal(expectedTransactionJournalUUID, transaction1.JournalUUID)
		is.Equal(expectedTransactionDescription, transaction1.Description)
		is.Equal(expectedTransactionMemo, transaction1.Memo)
		is.Equal(2, len(transaction1.TransactionLineItems))

		transactionLineItem1 := transaction1.TransactionLineItems[0]
		is.Equal(expectedLineItemUUID, transactionLineItem1.TransactionLineItemUUID)
		is.Equal(expectedLineItemTransactionUUID, transactionLineItem1.TransactionUUID)
		is.Equal(expectedLineItemDate, transactionLineItem1.Date)
		is.Equal(expectedLineItemAmount, transactionLineItem1.Amount)
		is.Equal(expectedLineItemAccount, transactionLineItem1.Account)
		is.Equal(expectedLineItemStatus, transactionLineItem1.Status)

		transactionLineItem2 := transaction1.TransactionLineItems[1]
		is.Equal(expectedLineItemUUID2, transactionLineItem2.TransactionLineItemUUID)
		is.Equal(expectedLineItemTransactionUUID2, transactionLineItem2.TransactionUUID)
		is.Equal(expectedLineItemDate2, transactionLineItem2.Date)
		is.Equal(expectedLineItemAmount2, transactionLineItem2.Amount)
		is.Equal(expectedLineItemAccount2, transactionLineItem2.Account)
		is.Equal(expectedLineItemStatus2, transactionLineItem2.Status)

		transaction2 := actualTransactions[1]
		is.Equal(expectedTransactionUUID2, transaction2.TransactionUUID)
		is.Equal(expectedTransactionJournalUUID2, transaction2.JournalUUID)
		is.Equal(expectedTransactionDescription2, transaction2.Description)
		is.Equal(expectedTransactionMemo2, transaction2.Memo)
		is.Equal(1, len(transaction2.TransactionLineItems))

		transactionLineItem3 := transaction2.TransactionLineItems[0]
		is.Equal(expectedLineItemUUID3, transactionLineItem3.TransactionLineItemUUID)
		is.Equal(expectedLineItemTransactionUUID3, transactionLineItem3.TransactionUUID)
		is.Equal(expectedLineItemDate3, transactionLineItem3.Date)
		is.Equal(expectedLineItemAmount3, transactionLineItem3.Amount)
		is.Equal(expectedLineItemAccount3, transactionLineItem3.Account)
		is.Equal(expectedLineItemStatus3, transactionLineItem3.Status)
	})
}
