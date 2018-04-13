package collector

import (
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/zamedic/multiPoolDashboard/balance"
)

type multiPoolHub struct {
}

func (s *multiPoolHub) getId() string {
	return "mph"
}

func (s *multiPoolHub) getName() string {
	return "Multi Pool Hub"
}

func newMultiPoolHub() pool {
	return &multiPoolHub{
	}
}

func (s *multiPoolHub) getBalance(apikey string) ([]balance.CoinValue, error) {
	response, err := http.Get(fmt.Sprintf("https://miningpoolhub.com/index.php?page=api&action=getuserallbalances&api_key=%v", apikey))
	if err != nil {
		return nil, err
	}
	data := &multiPoolBalanceResponse{}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	r := []balance.CoinValue{}
	for _, item := range data.Getuserallbalances.Data {
		r = append(r, balance.CoinValue{
			Name:  item.Coin,
			Coins: item.Ae_confirmed + item.Ae_unconfirmed + item.Confirmed + item.Exchange + item.Unconfirmed,
		})
	}
	return r, nil
}

type multiPoolBalanceResponse struct {
	Getuserallbalances multiPoolBalanceBody
}

type multiPoolBalanceBody struct {
	Version string
	Runtime float64
	Data    []multiPoolBalanceCoin
}

type multiPoolBalanceCoin struct {
	Coin           string
	Confirmed      float64
	Unconfirmed    float64
	Ae_confirmed   float64
	Ae_unconfirmed float64
	Exchange       float64
}
