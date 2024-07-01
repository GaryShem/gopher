package accrual

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/logging"
	"github.com/GaryShem/gopher/cmd/gophermart/internal/server/storage/repository"
)

var ErrOrderNotRegistered = errors.New("order not registered")
var ErrInternalAccrualError = errors.New("internal accrual error")
var ErrUnexpectedAccrualError = errors.New("unexpected accrual error")

type BonusTracker struct {
	Address     string
	restyClient *resty.Client
}

func NewBonusTracker(address string) *BonusTracker {
	logging.Log.Infoln("will try to address accrual at:", address)
	return &BonusTracker{
		Address:     address,
		restyClient: resty.New(),
	}
}

func (b *BonusTracker) UpdateOrder(orderID string) (*repository.Order, error) {
	return b.requestOrder(orderID, time.Second)
}

func (b *BonusTracker) requestOrder(orderID string, retryDelay time.Duration) (*repository.Order, error) {
	for {
		url := fmt.Sprintf("%s/api/orders/%s", b.Address, orderID)
		logging.Log.Infoln("trying to request accrual at url:", url)
		result := new(repository.Order)
		request := b.restyClient.R()
		response, err := request.Get(url)
		if err != nil {
			return nil, err
		}
		switch response.StatusCode() {
		case http.StatusOK:
			body := response.Body()
			if err = json.Unmarshal(body, result); err != nil {
				return result, err
			}
			result.Number = orderID
			if result.Status == "INVALID" || result.Status == "PROCESSED" {
				return result, nil
			}
			time.Sleep(retryDelay)
		case http.StatusNoContent:
			return nil, ErrOrderNotRegistered
		case http.StatusTooManyRequests:
			time.Sleep(retryDelay)
			continue
		case http.StatusInternalServerError:
			return nil, ErrInternalAccrualError
		default:
			return nil, ErrUnexpectedAccrualError
		}
	}
}
