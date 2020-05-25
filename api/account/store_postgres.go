package account

import (
	"database/sql"
)

// Postgres is the postgres accounts `Store` implementation.
type Postgres struct {
	DB *sql.DB
}

// SignUp persists a new account, user and membership of that user to the created account.
func (store *Postgres) SignUp(email string, password string) (accountID int64, userID int64, err error) {
	err = store.DB.QueryRow(signUpStatement, email, email, password).Scan(&accountID, &userID) // using email as account name
	return accountID, userID, err
}

// SignIn checks for a user matching the passed credentials
func (store *Postgres) SignIn(email string, password string) (userID int64, accounts map[int64]string, err error) {
	rows, err := store.DB.Query(signInStatement, email, password)
	if err != nil {
		e(err.Error())
		return -1, map[int64]string{}, err
	}

	ok := false
	for rows.Next() {
		ok = true
		accounts = map[int64]string{}
		var accountID int64
		var accountName string
		err = rows.Scan(&userID, &accountID, &accountName)
		if err != nil {
			e(err.Error())
			return -1, map[int64]string{}, err
		}
		accounts[accountID] = accountName
	}

	if !ok {
		err = sql.ErrNoRows
	}

	return userID, accounts, err
}

const signInStatement = `
SELECT users.id, memberships.account_id, accounts.name
FROM kenza.users 
INNER JOIN kenza.memberships ON users.id = memberships.user_id 
INNER JOIN kenza.accounts ON accounts.id = memberships.account_id                                                                                                                                                         
WHERE users.email = $1
AND password = crypt($2, password);
`

const signUpStatement = `
WITH insert_account AS (
		INSERT INTO kenza.accounts (email, name)
		VALUES ($1, $2)
		RETURNING email, id AS account_id
	), insert_user AS (
		INSERT INTO kenza.users (email, username, password)
		SELECT email, email, crypt($3, gen_salt('bf'))
		FROM insert_account
		RETURNING id AS user_id
	), insert_membership AS (
		INSERT INTO kenza.memberships (account_id, user_id)
		SELECT account_id, user_id		
		FROM insert_account, insert_user 
		RETURNING account_id, user_id
	)
SELECT account_id, user_id
FROM insert_account, insert_user`
