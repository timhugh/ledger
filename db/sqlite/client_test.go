package sqlite

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/matryer/is"
	_ "github.com/mattn/go-sqlite3"
	"github.com/timhugh/ledger"
	"testing"
)

var expectedJournalUUID = uuid.NewString()

const expectedJournalName = "Test Journal"

var expectedTransactionUUID = uuid.NewString()

const expectedTransactionDescription = "Test Transaction"
const expectedTransactionMemo = "Test Memo"

var expectedLineItemUUID = uuid.NewString()

const expectedLineItemDate = "2024-08-01"
const expectedLineItemAmount = 100
const expectedLineItemAccount = "Expenses:Test"
const expectedLineItemStatus = "pending"

func testClient(t *testing.T) (*Client, *sql.DB) {
	is := is.New(t)

	db, err := sql.Open("sqlite3", ":memory:")
	is.NoErr(err)

	client := &Client{db: db}
	is.NoErr(client.Migrate())

	return client, db
}

func TestCreateJournal(t *testing.T) {
	is := is.New(t)
	t.Run("creates journals", func(t *testing.T) {
		client, db := testClient(t)

		newJournal := &ledger.Journal{UUID: expectedJournalUUID, Name: expectedJournalName}
		err := client.CreateJournal(newJournal)
		is.NoErr(err)

		var actualUUID string
		var actualName string
		err = db.QueryRow("SELECT uuid, name FROM journals").Scan(&actualUUID, &actualName)
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
		is.NoErr(db.QueryRow("SELECT uuid, name FROM journals").Scan(&journalUUID, &journalName))

		if journalUUID == "" {
			t.Error("got empty uuid, want non-empty")
		}
		is.Equal(expectedJournalName, journalName)
	})
	t.Run("does not allow duplicate uuids", func(t *testing.T) {
		client, _ := testClient(t)

		existingJournal := &ledger.Journal{UUID: expectedJournalUUID, Name: expectedJournalName}
		is.NoErr(client.CreateJournal(existingJournal))

		newJournal := &ledger.Journal{UUID: expectedJournalUUID, Name: expectedJournalName}
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
	t.Run("returns a journal by id", func(t *testing.T) {
		client, db := testClient(t)

		_, err := db.Exec("INSERT INTO journals (uuid, name) VALUES (?, ?)", expectedJournalUUID, expectedJournalName)
		is.NoErr(err)

		actualJournal, err := client.GetJournal(expectedJournalUUID)
		is.NoErr(err)
		is.Equal(expectedJournalUUID, actualJournal.UUID)
		is.Equal(expectedJournalName, actualJournal.Name)
	})
	t.Run("returns a journal with transactions and line items", func(t *testing.T) {
		client, db := testClient(t)

		_, err := db.Exec("INSERT INTO journals (uuid, name) VALUES (?, ?)", expectedJournalUUID, expectedJournalName)
		is.NoErr(err)
		_, err = db.Exec("INSERT INTO transactions (uuid, journal_uuid, description, memo) VALUES (?, ?, ?, ?)", expectedTransactionUUID, expectedJournalUUID, expectedTransactionDescription, expectedTransactionMemo)
		is.NoErr(err)
		_, err = db.Exec("INSERT INTO transaction_line_items (uuid, transaction_uuid, date, amount, account, status) VALUES (?, ?, ?, ?, ?, ?)", expectedLineItemUUID, expectedTransactionUUID, expectedLineItemDate, expectedLineItemAmount, expectedLineItemAccount, expectedLineItemStatus)
		is.NoErr(err)

		actualJournal, err := client.GetJournal(expectedJournalUUID)
		is.NoErr(err)
		is.Equal(expectedJournalUUID, actualJournal.UUID)
		is.Equal(expectedJournalName, actualJournal.Name)
		is.Equal(1, len(actualJournal.Transactions))

		transaction := actualJournal.Transactions[0]
		is.Equal(expectedTransactionUUID, transaction.UUID)
		is.Equal(expectedTransactionDescription, transaction.Description)
		is.Equal(expectedTransactionMemo, transaction.Memo)
		is.Equal(1, len(transaction.LineItems))

		lineItem := transaction.LineItems[0]
		is.Equal(expectedLineItemUUID, lineItem.UUID)
		is.Equal(expectedLineItemDate, lineItem.Date)
		is.Equal(expectedLineItemAmount, lineItem.Amount)
		is.Equal(expectedLineItemAccount, lineItem.Account)
		is.Equal(expectedLineItemStatus, string(lineItem.Status))
	})
}
