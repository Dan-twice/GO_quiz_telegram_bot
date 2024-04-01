package database

import "github.com/jmoiron/sqlx"

var (
	db     *sqlx.DB
	schema = `
DROP TABLE dictionary;

CREATE TABLE dictionary (
    id integer,
    word text,
    description text,
    context text
);
`
)

type Dictionary struct {
	Word        string
	Description string
	Context     string
}
