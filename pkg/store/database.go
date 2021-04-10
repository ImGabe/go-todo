package store

var CreateDatabaseStatement = `
	CREATE TABLE IF NOT EXISTS task (
		id   INTEGER NOT NULL PRIMARY KEY,
		description TEXT    NOT NULL,
		done BOOL    NOT NULL
	)
`
