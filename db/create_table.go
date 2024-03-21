package db

func CreateTables() {

	userTable := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT,
			transaction_pin: TEXT NOT NULL,
			tag TEXT UNIQUE NOT NULL,
			phone TEXT UNIQUE NOT NULL,
			is_verified BOOLEAN NOT NULL,
			created_at TEXT,
			updated_at TEXT,
			deleted_at TEXT
		)
	`
	_, err := MainDB.Exec(userTable)
	if err != nil {
		panic(err)
	}

	walletTable := `
		CREATE TABLE IF NOT EXISTS wallets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			balance INTEGER NOT NULL,
			created_at TEXT,
			updated_at TEXT,
			deleted_at TEXT,
			FOREIGN KEY (user_id) REFERENCES users(id)
			`
	_, err = MainDB.Exec(walletTable)
	if err != nil {
		panic(err)
	}

}
