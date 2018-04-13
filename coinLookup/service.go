package coinLookup

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type Service interface {
	GetSymbolFromName(string) string
	GetNameFromSymbol(string) string
}

type coinlookup struct {
	coinMapSymbol map[string]string
	coinMapName map[string]string
}

func (s *coinlookup) GetSymbolFromName(in string) string {
	return s.coinMapName[in]
}

func (s *coinlookup) GetNameFromSymbol(in string) string {
	return s.coinMapSymbol[in]
}

func NewService() Service{
	s := &coinlookup{}
	response, err := http.Get("https://github.com/crypti/cryptocurrencies/blob/master/cryptocurrencies.json")
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	s.coinMapSymbol= make(map[string]string)
	err = json.Unmarshal(b, &s.coinMapSymbol)
	if err != nil {
		panic(err)
	}

	s.coinMapName = make(map[string]string)
	for key, value := range s.coinMapSymbol {
		s.coinMapName[value]  = key
	}
	return s
}




