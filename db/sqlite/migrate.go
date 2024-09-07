package sqlite

import (
	"context"
	"github.com/timhugh/ctxlogger"
)

func (c *Client) Migrate(ctx context.Context) error {
	for _, migration := range migrations {
		// TODO: track migrations in table
		ctxlogger.Info(ctx, "running migration %d: %s", migration.ID, migration.Name)
		_, err := c.db.ExecContext(ctx, migration.SQL)
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
            memo text not null,
            FOREIGN KEY (journal_uuid) REFERENCES journals(journal_uuid)
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
            status text not null default 'pending',
            FOREIGN KEY (transaction_uuid) REFERENCES transactions(transaction_uuid)
        );`,
	},
	{
		4, "create users table",
		`create table users (
            user_uuid string primary key not null,
            login text not null,
            password_hash text not null,
            password_salt text not null
        );`,
	},
	{
		5, "create sessions table",
		`create table sessions (
            session_uuid string primary key not null,
            user_uuid string not null,
            FOREIGN KEY (user_uuid) REFERENCES users(user_uuid)
        );`,
	},
}
