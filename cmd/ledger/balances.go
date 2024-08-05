package main

import (
    "github.com/urfave/cli/v2"
    "fmt"
    "os"
    "text/tabwriter"
)

var balancesCmd = cli.Command{
    Name: "balances",
    Aliases: []string{"bal", "b"},
    Category: "report",
    Usage: "print all account balances",
    Action: func(c *cli.Context) error {
        l := loadLedger()
        balances := l.Balances()
        w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
        for account, balance := range balances {
            fmt.Fprintf(w, "%s\t%d\n", account, balance)
        }
        w.Flush()
        return nil
    },
}

