package db

func CreateTables() {

	userTable := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			transaction_pin: TEXT NOT NULL,
			tag TEXT UNIQUE NOT NULL,
			phone TEXT UNIQUE NOT NULL,
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
