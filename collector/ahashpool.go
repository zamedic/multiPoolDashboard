package collector

import (
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/zamedic/multiPoolDashboard/balance"
)

type ahashpool struct {
}

func (s *ahashpool) getId() string {
	return "ahashpool"
}

func (s *ahashpool) getName() string {
	return "ahashpool"
}

func newAhashPool() pool {
	return &ahashpool{}
}

func (s *ahashpool) getBalance(apikey string) ([]balance.CoinValue, error) {
	response, err := http.Get(fmt.Sprintf("http://www.ahashpool.com/api/wallet?address=%v", apikey))
	if err != nil {
		return nil, err
	}

	data := ahashpoolBalance{}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	r := balance.CoinValue{Name: data.Currency, Coins: data.Total_unpaid}
	return []balance.CoinValue{r}, nil
}

type ahashpoolBalance struct {
	Currency     string
	Unsold       float64
	Balance      float64
	Total_unpaid float64
	Total_paid   float64
	Total_earned float64
}
