package application

import "errors"

var (
	ErrInvalidBody             = errors.New("Invalid body")
	ErrInvalidUnmarshal        = errors.New("Invalid unmarshal")
	ErrValidateStruct          = errors.New("Invalid validate struct")
	ErrInvalidGeneratePassHash = errors.New("Invalid generate password hash")
	ErrUserAlreadyExists       = errors.New("User already exists")
	ErrUserNotFound            = errors.New("User not found")
	ErrInvalidGetDataFromDB    = errors.New("Invalid get data from DB")

	ErrInvalidUsernameOrPassword = errors.New("Invalid username or password")
	ErrInvalidAddUser            = errors.New("Invalid add user")
	ErrInvalidGenerateToken      = errors.New("Invalid generate token")
	ErrInvalidUpdateData         = errors.New("Invalid update data")

	ErrInvalidInterval = errors.New("Invalid interval")

	ErrInvalidAddBunker  = errors.New("Invalid add bunker")
	ErrAddBunker         = errors.New("Can not add bunker")
	ErrBunkerNotFound    = errors.New("Bunker not found")
	ErrInvalidDeleteData = errors.New("Invalid delete data")
)
