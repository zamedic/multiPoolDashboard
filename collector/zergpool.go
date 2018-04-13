package collector

import (
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/zamedic/multiPoolDashboard/balance"
)

type zergpool struct {
}

func newZergPool() pool{
	return &zergpool{}
}

func (s zergpool) getBalance(apikey string) ([]balance.CoinValue, error) {
	response, err := http.Get(fmt.Sprintf("http://api.zergpool.com:8080/api/wallet?address=%v", apikey))
	if err != nil {
		return nil, err
	}
	z := zerg{}
	err = json.NewDecoder(response.Body).Decode(&z)
	if err != nil {
		return nil, err
	}
	return []balance.CoinValue{{Name:z.Currency,Coins:z.Unpaid}}, nil

}

func (s zergpool) getId() string {
	return "zerg"
}

func (s zergpool) getName() string {
	return "zergpool"
}

type zerg struct {
	Currency string
	Unsold   float64
	Balance  float64
	Unpaid   float64
	Paid24h  float64
	total    float64
}
