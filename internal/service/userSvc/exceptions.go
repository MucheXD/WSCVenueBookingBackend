package userSvc

import (
	"errors"
)

var (
	ErrRegisterInfoInvalid                 = errors.New("one of the registration information is invalid")
	ErrLoginTokenExpired                   = errors.New("login token has expired")
	ErrInvalidLoginTokenLength             = errors.New("invalid login token length")
	ErrLoginSessionIDInvalid               = errors.New("login session ID is invalid")
	ErrParseValidBeforeTimestamp           = errors.New("failed to parse valid-before timestamp")
	ErrLoginTokenSignatureMismatch         = errors.New("login token signature mismatch")
	ErrDecodeLoginTokenPasswordHash        = errors.New("failed to decode login token password hash")
	ErrLoginTokenSaltValidBeforeOutOfRange = errors.New("generated login token salt valid before is out of range")
	ErrUserPasswordHashInvalid             = errors.New("user password hash in database is invalid")
	ErrUserNotVerified                     = errors.New("user is not verified for changing password")
	ErrNewPasswordOrSaltInvalid            = errors.New("new password or salt is invalid")

	ErrLoginTokenSaltSecretNotConfigured = errors.New("login token salt secret is not configured")
	ErrCheckLoginSessionIDInRedis        = errors.New("failed to check login session ID in redis")
	ErrStoreLoginSessionIDInRedis        = errors.New("failed to store login session ID in redis")

	ErrUsernameAlreadyExists = errors.New("username already exists")

	ErrCheckUsernameExistsInDB = errors.New("failed to check if username exists in database")
	ErrCreateUserInDB          = errors.New("failed to create user in database")
	ErrUpdateUserInDB          = errors.New("failed to edit user in database")
)
