insert into journals (journal_uuid, name)
VALUES ('234e5632-8864-46fe-a1a9-a95b4d03b147', '2024 Journal');

begin;
insert into transactions (transaction_uuid, journal_uuid, description, memo)
VALUES ('b7e6f708-deb0-470f-acbd-42220556f66c', '234e5632-8864-46fe-a1a9-a95b4d03b147',
        'Starting balances', '');
insert into transaction_line_items (transaction_line_item_uuid, transaction_uuid, date, amount, account)
VALUES ('3d54691e-0cfa-4aaa-b0b9-a5b6965f5113', 'b7e6f708-deb0-470f-acbd-42220556f66c',
        '2024-01-01', 1000, 'Assets:Checking'),
       ('43256c9c-f9ad-4107-9dbf-4c2303f700e1', 'b7e6f708-deb0-470f-acbd-42220556f66c',
        '2024-01-01', -1000, 'Equity:Opening Balances');
commit;

begin;
insert into transactions (transaction_uuid, journal_uuid, description, memo)
VALUES ('298a3ea4-ce78-47b9-976b-d743934a9690', '234e5632-8864-46fe-a1a9-a95b4d03b147', 'Rent',
        'External Withrawal - YOUR LANDLORD');
insert into transaction_line_items (transaction_line_item_uuid, transaction_uuid, date, amount, account)
VALUES ('5127d7a2-421b-4805-bc23-79eb859ac470', '298a3ea4-ce78-47b9-976b-d743934a9690',
        '2024-01-01', 500, 'Expenses:Rent'),
       ('bce95758-b558-47c9-94a3-172f3fae1032', '298a3ea4-ce78-47b9-976b-d743934a9690',
        '2024-01-01', -500, 'Assets:Checking');
commit;
