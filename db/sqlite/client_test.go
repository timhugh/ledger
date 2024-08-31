package sqlite

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/timhugh/ledger"
    "github.com/stretchr/testify/assert"
	"testing"
)

var expectedJournal = &ledger.Journal{
    JournalUUID: "234e5632-8864-46fe-a1a9-a95b4d03b147",
    Name: "Test Journal",
    Transactions: []*ledger.Transaction{
        {
            TransactionUUID: "7115904B-0DA6-46C5-B385-0276095666B0",
            JournalUUID: "234e5632-8864-46fe-a1a9-a95b4d03b147",
            Description: "Test Transaction",
            Memo: "Test Memo",
            TransactionLineItems: []*ledger.TransactionLineItem{
                {
                    TransactionLineItemUUID: "A3C2DE2D-D959-40D1-9CFC-878C66D15204",
                    TransactionUUID: "7115904B-0DA6-46C5-B385-0276095666B0",
                    Date: "2024-08-01",
                    Amount: 100,
                    Account: "Expenses:Test",
                    Status: ledger.Cleared,
                },
                {
                    TransactionLineItemUUID: "A3C2DE2D-D959-40D1-9CFC-878C66D15205",
                    TransactionUUID: "7115904B-0DA6-46C5-B385-0276095666B0",
                    Date: "2024-08-02",
                    Amount: 200,
                    Account: "Expenses:Test2",
                    Status: ledger.Cleared,
                },
            },
        },
        {
            TransactionUUID: "7115904B-0DA6-46C5-B385-0276095666B1",
            JournalUUID: "234e5632-8864-46fe-a1a9-a95b4d03b147",
            Description: "Test Transaction 2",
            Memo: "Test Memo 2",
            TransactionLineItems: []*ledger.TransactionLineItem{
                {
                    TransactionLineItemUUID: "A3C2DE2D-D959-40D1-9CFC-878C66D15206",
                    TransactionUUID: "7115904B-0DA6-46C5-B385-0276095666B1",
                    Date: "2024-08-03",
                    Amount: 300,
                    Account: "Expenses:Test3",
                    Status: ledger.Pending,
                },
            },
        },
    },
}

func insertExpectedJournal(assert *assert.Assertions, db *sqlx.DB) {
		_, err := db.Exec("INSERT INTO journals (journal_uuid, name) VALUES (?, ?)", expectedJournal.JournalUUID, expectedJournal.Name)
		assert.NoError(err)
        expectedTransaction1 := expectedJournal.Transactions[0]
		_, err = db.Exec("INSERT INTO transactions (transaction_uuid, journal_uuid, description, memo) VALUES (?, ?, ?, ?)", expectedTransaction1.TransactionUUID, expectedTransaction1.JournalUUID, expectedTransaction1.Description, expectedTransaction1.Memo)
		assert.NoError(err)
        expectedLineItem1 := expectedTransaction1.TransactionLineItems[0]
		_, err = db.Exec("INSERT INTO transaction_line_items (transaction_line_item_uuid, transaction_uuid, date, amount, account, status) VALUES (?, ?, ?, ?, ?, ?)", expectedLineItem1.TransactionLineItemUUID, expectedLineItem1.TransactionUUID, expectedLineItem1.Date, expectedLineItem1.Amount, expectedLineItem1.Account, expectedLineItem1.Status)
		assert.NoError(err)
        expectedLineItem2 := expectedTransaction1.TransactionLineItems[1]
		_, err = db.Exec("INSERT INTO transaction_line_items (transaction_line_item_uuid, transaction_uuid, date, amount, account, status) VALUES (?, ?, ?, ?, ?, ?)", expectedLineItem2.TransactionLineItemUUID, expectedLineItem2.TransactionUUID, expectedLineItem2.Date, expectedLineItem2.Amount, expectedLineItem2.Account, expectedLineItem2.Status)
		assert.NoError(err)
        expectedTransaction2 := expectedJournal.Transactions[1]
		_, err = db.Exec("INSERT INTO transactions (transaction_uuid, journal_uuid, description, memo) VALUES (?, ?, ?, ?)", expectedTransaction2.TransactionUUID, expectedTransaction2.JournalUUID, expectedTransaction2.Description, expectedTransaction2.Memo)
		assert.NoError(err)
        expectedLineItem3 := expectedTransaction2.TransactionLineItems[0]
		_, err = db.Exec("INSERT INTO transaction_line_items (transaction_line_item_uuid, transaction_uuid, date, amount, account, status) VALUES (?, ?, ?, ?, ?, ?)", expectedLineItem3.TransactionLineItemUUID, expectedLineItem3.TransactionUUID, expectedLineItem3.Date, expectedLineItem3.Amount, expectedLineItem3.Account, expectedLineItem3.Status)
		assert.NoError(err)
}

func testClient(t *testing.T) (*assert.Assertions, *Client, *sqlx.DB) {
    assert := assert.New(t)

	db, err := sqlx.Open("sqlite3", ":memory:")
    assert.NoError(err)

	client := &Client{db: db}
    assert.NoError(client.Migrate())

	return assert, client, db
}

func TestCreateJournal(t *testing.T) {
	t.Run("satisfies JournalCreator interface", func(t *testing.T) {
		assert, client, _ := testClient(t)
		_, ok := interface{}(client).(ledger.JournalCreator)
        assert.True(ok)
	})
	t.Run("creates journals", func(t *testing.T) {
		assert, client, db := testClient(t)

        assert.NoError(client.CreateJournal(expectedJournal))

		var actualUUID string
		var actualName string
		assert.NoError(db.QueryRow("SELECT journal_uuid, name FROM journals").Scan(&actualUUID, &actualName))

        assert.Equal(expectedJournal.JournalUUID, actualUUID)
        assert.Equal(expectedJournal.Name, actualName)
	})
	t.Run("automatically assigns IDs", func(t *testing.T) {
		assert, client, db := testClient(t)

		newJournal := &ledger.Journal{Name: expectedJournal.Name}
        assert.NoError(client.CreateJournal(newJournal))

		var journalUUID string
		var journalName string
        assert.NoError(db.QueryRow("SELECT journal_uuid, name FROM journals").Scan(&journalUUID, &journalName))

        assert.NotEmpty(journalUUID)
        assert.Equal(expectedJournal.Name, journalName)
	})
	t.Run("does not allow duplicate uuids", func(t *testing.T) {
		assert, client, _ := testClient(t)

        assert.NoError(client.CreateJournal(expectedJournal))
        assert.Error(client.CreateJournal(expectedJournal))
	})
	t.Run("requires a name", func(t *testing.T) {
		// skipping for the moment -- go is passing the zero value of a string "" instead of nil
		// so the not null constraint is being satisfied
		t.Skip("not implemented")
		assert, client, _ := testClient(t)

		actualJournal := &ledger.Journal{}
		err := client.CreateJournal(actualJournal)
        assert.Error(err)
	})
}
func TestGetJournal(t *testing.T) {
	t.Run("satisfies JournalGetter interface", func(t *testing.T) {
		assert, client, _ := testClient(t)
		_, ok := interface{}(client).(ledger.JournalGetter)
        assert.True(ok)
	})
	t.Run("returns a journal by id", func(t *testing.T) {
		assert, client, db := testClient(t)

        insertExpectedJournal(assert, db)

		actualJournal, err := client.GetJournal(expectedJournal.JournalUUID)
        assert.NoError(err)
        assert.Equal(expectedJournal.JournalUUID, actualJournal.JournalUUID)
        assert.Equal(expectedJournal.Name, actualJournal.Name)
	})
    t.Run("returns an error if the journal does not exist", func(t *testing.T) {
        assert, client, _ := testClient(t)
        journal, err := client.GetJournal("nonexistent")
        assert.Equal(NoRecordError, err)
        assert.Nil(journal)
    })
	t.Run("returns a journal with transactions and line items", func(t *testing.T) {
		assert, client, db := testClient(t)

        insertExpectedJournal(assert, db)

		actualJournal, err := client.GetJournal(expectedJournal.JournalUUID)
		assert.NoError(err)
        assert.EqualValues(expectedJournal, actualJournal)
	})
}
func TestGetTransactions(t *testing.T) {
	t.Run("satisfies TransactionLister interface", func(t *testing.T) {
		assert, client, _ := testClient(t)
		_, ok := interface{}(client).(ledger.TransactionLister)
        assert.True(ok)
	})
	t.Run("returns transactions for a journal", func(t *testing.T) {
		assert, client, db := testClient(t)

        insertExpectedJournal(assert, db)

		actualTransactions, err := client.GetTransactions(expectedJournal.JournalUUID)
		assert.NoError(err)

        assert.EqualValues(expectedJournal.Transactions, actualTransactions)
    })
}
