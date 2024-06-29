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

type BonusTracker struct {
	Address string
}

func NewBonusTracker(address string) *BonusTracker {
	logging.Log.Infoln("will try to address accrual at:", address)
	return &BonusTracker{Address: address}
}

func (b *BonusTracker) UpdateOrder(orderID string) (*repository.Order, error) {
	client := resty.New()
	url := fmt.Sprintf("%s/api/orders/%s", b.Address, orderID)
	logging.Log.Infoln("trying to request accrual at url:", url)
	request := client.R()
	result := new(repository.Order)
	for {
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
			time.Sleep(time.Millisecond * 100)
		case http.StatusNoContent:
			return nil, ErrOrderNotRegistered
		case http.StatusTooManyRequests:
			time.Sleep(time.Second * 1)
			continue
		case http.StatusInternalServerError:
			return nil, ErrInternalAccrualError
		}
	}
}
