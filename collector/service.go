package collector

import (
	"github.com/zamedic/multiPoolDashboard/user"
	"time"
	"log"
	"github.com/zamedic/multiPoolDashboard/balance"
)

type pool interface {
	getBalance(apikey string) ([]balance.CoinValue, error)
	getId() string
	getName() string
}

type service struct {
	services []pool
	userStore    user.Store
	balanceStore balance.Store
}

func NewCollector(userStore user.Store, balanceStore balance.Store) {
	s := service{userStore:userStore,balanceStore:balanceStore}
	s.services = []pool{
		newMultiPoolHub(),
		newAhashPool(),
		newZergPool(),
	}

	go func() {
		s.runScan()
	}()
}

func (s *service) runScan() {
	for {
		s.userStore.ScanUsers(s.iterateEmails)
		time.Sleep(10 * time.Minute)
	}
}

func (s *service) iterateEmails(emails []string) {
	log.Println("Starting email scan")
	for _, email := range emails {
		log.Printf("Scanning email: %v", email)
		coins := map[string][]balance.CoinValue{}
		for _, pool := range s.services {
			log.Printf("Scanning pool: %v", pool.getId())
			id, _ := s.userStore.GetKey(email, pool.getId())
			if len(id) > 0 {
				log.Printf("Looking up ID %v", id)
				c, _ := pool.getBalance(id)
				coins[pool.getId()] = c
			}
		}
		if len(coins) > 0 {
			s.balanceStore.SaveBalance(email, coins)
		}
	}
}
