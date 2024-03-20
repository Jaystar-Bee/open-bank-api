package db

func CreateTables() {

	userTable := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			email TEXT NOT NULL,
			password TEXT NOT NULL,
			transaction_pin: TEXT NOT NULL,
			account_number TEXT NOT NULL,
			phone TEXT NOT NULL,
			created_at TEXT,
			updated_at TEXT,
			deleted_at TEXT
		)
	`
	_, err := MainDB.Exec(userTable)
	if err != nil {
		panic(err)
	}

}
