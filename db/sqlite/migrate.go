package sqlite

import (
	"log"
)

func (c *Client) Migrate() error {
	for _, migration := range migrations {
		// TODO: track migrations in table
		log.Printf("running migration %d: %s\n", migration.ID, migration.Name)
		_, err := c.db.Exec(migration.SQL)
		if err != nil {
			return err
		}
	}
	return nil
}

type Migration struct {
	ID   int
	Name string
	SQL  string
}

var migrations = []Migration{
	{
		1, "create journals table",
		`create table journals (
            journal_uuid string primary key not null,
            name text not null
        );`,
	},
	{
		2, "create transactions table",
		`create table transactions (
            transaction_uuid string primary key not null,
            journal_uuid int not null,
            description text,
            memo text not null
        );`,
	},
	{
		3, "create line_items table",
		`create table transaction_line_items (
            transaction_line_item_uuid string primary key not null,
            transaction_uuid int not null,
            date text not null,
            amount int not null,
            account text not null,
            status text not null default 'pending'
        );`,
	},
}
