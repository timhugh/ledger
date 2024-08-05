package ledger

type Transaction struct {
    LineItems []LineItem
}

func (t *Transaction) Cleared() bool {
    for _, lineItem := range t.LineItems {
        if lineItem.Status != Cleared {
            return false
        }
    }
    return true
}

func (t *Transaction) Date() Date {
    lowestDate := t.LineItems[0].Date
    for _, lineItem := range t.LineItems {
        if lineItem.Date < lowestDate {
            lowestDate = lineItem.Date
        }
    }
    return lowestDate
}

func (t *Transaction) Valid() bool {
    var sum int
    for _, lineItem := range t.LineItems {
        sum += int(lineItem.Amount)
    }
    return sum == 0
}

type LineItem struct {
    Account Account
    Amount Amount
    Date Date
    Status Status
}

type Account string
type Date string
type Amount int

type Status string
const (
    Cleared Status = "cleared"
    Pending Status = "pending"
)

