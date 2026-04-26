package ws

import (
	"errors"
)

var (
	ErrTokenExpired     = errors.New("Token expired")
	ErrTokenNotValidYet = errors.New("Token not valid yet")
	ErrInvalidAudience  = errors.New("Invalid audience")
	ErrInvalidIssuer    = errors.New("Invalid issuer")
	ErrInvalidToken     = errors.New("Invalid token")
	ErrTokenError       = errors.New("Token error")
	ErrInvalidSubject   = errors.New("Invalid subject")
)
