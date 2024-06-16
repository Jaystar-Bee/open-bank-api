package db

func CreateTables() {

	userTable := `
		CREATE TABLE users (
			id SERIAL PRIMARY KEY,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			transaction_pin TEXT NOT NULL,
			tag TEXT UNIQUE NOT NULL,
			phone TEXT,
			avatar TEXT,
			account_is_deactivated BOOLEAN NOT NULL,
			is_verified BOOLEAN NOT NULL,
			created_at TIMESTAMP,
			updated_at TIMESTAMP,
			deleted_at TIMESTAMP
		)
	`
	_, err := MainDB.Exec(userTable)
	if err != nil {
		panic(err)
	}

	walletTable := `
		CREATE TABLE wallets(
			id SERIAL PRIMARY KEY,
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
		CREATE TABLE transactions(
			id SERIAL PRIMARY KEY,
			sender INTEGER NOT NULL,
			sender_wallet INTEGER NOT NULL,
			receiver INTEGER NOT NULL,
			receiver_wallet INTEGER NOT NULL,
			amount FLOAT NOT NULL,
			status TEXT NOT NULL,
			remarks TEXT,
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

	request_table := `
		CREATE TABLE requests(
			id SERIAL PRIMARY KEY,
			requester INTEGER NOT NULL,
			giver INTEGER NOT NULL,
			amount FLOAT NOT NULL,
			status TEXT NOT NULL,
			remarks TEXT,
			created_at TEXT NOT NULL,
			updated_at TEXT,
			deleted_at TEXT,
			FOREIGN KEY (requester) REFERENCES users(id),
			FOREIGN KEY (giver) REFERENCES users(id)
			)
		`
	_, err = MainDB.Exec(request_table)
	if err != nil {
		panic(err)
	}

}
