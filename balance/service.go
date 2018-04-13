package balance

import "time"

type Service interface {
	getBalanceByPool(email string) []balance
	getBalanceByCoin(email, coin string) map[int64]float64
}

func NewService(store Store) Service{
	return &service{store:store}
}

type service struct {
	store Store
}

func (s service) getBalanceByCoin(email, coin string) map[int64]float64 {
	t := time.Now().Add(time.Hour * time.Duration(-12))
	return s.store.getBalanceByCoinName(email,coin, t)
}

func (s service) getBalanceByPool(email string) []balance {
	t := time.Now().Add(time.Hour * time.Duration(-12))
	result := s.store.getBalanceByPool(email, t)
	var r []balance
	for timestamp, pool := range result {
		r = append(r, balance{Timestamp: timestamp, Pools: pool})
	}
	return r
}

type balance struct {
	Timestamp int64
	Pools     map[string][]CoinValue
}
