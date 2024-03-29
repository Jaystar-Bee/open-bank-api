package db

func CreateTables() {

	userTable := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			transaction_pin TEXT NOT NULL,
			tag TEXT UNIQUE NOT NULL,
			phone TEXT,
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
		CREATE TABLE IF NOT EXISTS wallets(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER UNIQUE NOT NULL,
			balance FLOAT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT,
			deleted_at TEXT,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`

	_, err = MainDB.Exec(walletTable)
	if err != nil {
		panic(err)
	}

	transaction_table := `
		CREATE TABLE IF NOT EXISTS transactions(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			sender INTEGER NOT NULL,
			sender_wallet INTEGER NOT NULL,
			receiver INTEGER NOT NULL,
			receiver_wallet INTEGER NOT NULL,
			amount FLOAT NOT NULL,
			status BOOLEAN NOT NULL,
			type TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT,
			deleted_at TEXT,
			FOREIGN KEY (sender) REFERENCES users(id),
			FOREIGN KEY (receiver) REFERENCES users(id),
			FOREIGN KEY (sender_wallet) REFERENCES wallets(id),
			FOREIGN KEY (receiver_wallet) REFERENCES wallets(id)
		)
	`

	_, err = MainDB.Exec(transaction_table)
	if err != nil {
		panic(err)
	}

}
