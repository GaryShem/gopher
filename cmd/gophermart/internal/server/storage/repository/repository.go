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

type User struct {
	ID       int
	Name     string
	Password string
}

type Order struct {
	Number     string `json:"number"`
	Status     string `json:"status"`
	Accrual    int    `json:"accrual,omitempty"`
	UploadedAt string `json:"uploaded_at"`
}

type BalanceInfo struct {
	Current   float64
	Withdrawn float64
}
type WithdrawalInfo struct {
	Order       string
	Sum         float64
	ProcessedAt string
}
type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Repository interface {
	UserRegister(name, password string) error
	UserLogin(name, password string) (int, error)
	GetUserByName(name string) (User, error)

	OrderUpload(userID int, orderID string) error
	OrderGet(userID int) ([]Order, error)
	GetOrdersByUserID(userID int) ([]Order, error)

	BalanceList(userID int) (BalanceInfo, error)
	BalanceWithdraw(userID int, orderID string, amount float64) error
	BalanceWithdrawInfo(userID int) ([]WithdrawalInfo, error)
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
