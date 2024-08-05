package ledger

type Ledger struct {
    Transactions []Transaction
}

func (l *Ledger) Balances() map[Account]Amount {
    balances := make(map[Account]Amount)
    for _, transaction := range l.Transactions {
        for _, lineItem := range transaction.LineItems {
            balances[lineItem.Account] += lineItem.Amount
        }
    }
    return balances
}

func (l* Ledger) Register() []Transaction {
    return l.Transactions
}

