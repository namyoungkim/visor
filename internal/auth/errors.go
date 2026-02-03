package auth

import "errors"

// ErrNoCredentials is returned when no credentials are available.
var ErrNoCredentials = errors.New("no credentials available")

// ErrExpiredCredentials is returned when credentials have expired.
var ErrExpiredCredentials = errors.New("credentials have expired")

// ErrKeychainAccess is returned when keychain access fails.
var ErrKeychainAccess = errors.New("keychain access denied")
