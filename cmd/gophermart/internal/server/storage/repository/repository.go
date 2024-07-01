package repository

import (
	"errors"
)

// user errors

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user not found")
var ErrInvalidPassword = errors.New("invalid password")

// order errors

var ErrOrderUploadedSameUser = errors.New("order already uploaded by same user")
var ErrOrderUploadedDifferentUser = errors.New("order already uploaded by different user")
var ErrOrderIDFormatInvalid = errors.New("order ID format is invalid")

// balance errors

var ErrBalanceNotEnough = errors.New("balance not enough")
var ErrNoWithdrawals = errors.New("no withdrawals for current user")

type User struct {
	ID       int
	Name     string
	Password string
}

type Order struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	UserID     int     `json:"-"`
	Accrual    float64 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at,omitempty"`
}

type BalanceInfo struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type WithdrawalInfo struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at,omitempty"`
}

type CredentialRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Repository interface {
	RegisterUser(name, password string) error
	LoginUser(name, password string) (int, error)
	GetUserByName(name string) (User, error)
	CheckUserCredentials(name, password string) (int, error)

	UploadOrder(userID int, orderID string) error
	GetOrdersByUser(userID int) ([]Order, error)
	ProcessOrderUpdate(userID int, orderID string) error

	ListBalance(userID int) (BalanceInfo, error)
	WithdrawBalance(userID int, orderID string, amount float64) error
	GetBalanceWithdrawInfo(userID int) ([]WithdrawalInfo, error)
}

func ValidateOrderID(orderID string) error {
	var multPosition int
	if len(orderID)%2 == 0 {
		multPosition = 0
	} else {
		multPosition = 1
	}

	fullSum := 0
	for i := range len(orderID) {
		digitSum := int(orderID[i] - '0')
		if i%2 == multPosition {
			digitSum *= 2
		}
		if digitSum >= 10 {
			digitSum = digitSum - 10 + 1
		}
		fullSum += digitSum
	}
	if fullSum%10 != 0 {
		return ErrOrderIDFormatInvalid
	}
	return nil
}
