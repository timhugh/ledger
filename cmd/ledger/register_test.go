package main_test

import (
    "github.com/timhugh/ledger"
    "github.com/timhugh/ledger/cmd/ledger"
)

func ExampleRegister() {
    ledger := ledger.Ledger{
        Transactions: []ledger.Transaction{
            {
                LineItems: []ledger.LineItem{
                    { Date: "2024-08-01", Account: "Assets:Checking", Amount: 1000 },
                    { Date: "2024-08-01", Account: "Assets:Savings", Amount: 10000 },
                    { Date: "2024-08-01", Account: "Expenses:Groceries", Amount: -200 },
                    { Date: "2024-08-01", Account: "Expenses:Rent", Amount: -1000 },
                },
            },
        },
    }
    main.Register(ledger)

    // Output:
    // Date          Account               Amount
    // 2024-08-01    Assets:Checking       1000
    // 2024-08-01    Assets:Savings        10000
    // 2024-08-01    Expenses:Groceries    -200
    // 2024-08-01    Expenses:Rent         -1000
}
