package account

import "errors"

var (
	// errBadCredentials is thrown when the email / password combination does not match a user record.
	errBadCredentials = errors.New("Username and / or password provided is wrong")

	// errEmailNotProvided is thrown when email is not provided.
	errEmailNotProvided = errors.New("No email was provided")

	// errAccountEmailAlreadyExists is thrown when the email provided is already associated with an account.
	errAccountEmailAlreadyExists = errors.New("Email already exists")

	// errInternalStoreError is an "unkown / generic" error equivalent in case we fail to catch a specific error type.
	errInternalStoreError = errors.New("Internal store error")
)
