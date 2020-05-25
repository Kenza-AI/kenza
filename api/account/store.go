package account

// Store - accounts store abstraction
type Store interface {
	SignUp(email string, password string) (accountID int64, userID int64, err error)
	SignIn(email string, password string) (userID int64, accountIDs map[int64]string, err error)
}
