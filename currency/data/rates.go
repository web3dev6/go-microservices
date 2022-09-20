package data

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/go-hclog"
)

type ExchangeRates struct {
	log   hclog.Logger
	rates map[string]float64
}

// GetRate fetches currency rate for given base & destination currencies
func (e *ExchangeRates) GetRate(base, dest string) (float64, error) {
	br, ok := e.rates[base]
	if !ok {
		return 0, fmt.Errorf("rate not found for currency %s", base)
	}
	dr, ok := e.rates[dest]
	if !ok {
		return 0, fmt.Errorf("rate not found for currency %s", dest)
	}
	return dr / br, nil
}

// NewRates instantiates a new ExchangeRates
func NewRates(l hclog.Logger) (*ExchangeRates, error) {
	er := &ExchangeRates{log: l, rates: map[string]float64{}}
	er.fetchRatesFromECB()
	return er, nil
}

// GetRates fetches currency rates against EUR for various currencies from eu-central-bank
func (e *ExchangeRates) fetchRatesFromECB() error {
	resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status_code 200, got %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	// resp in XML : set rates map in ExchangeRates instance
	md := &Cubes{}
	xml.NewDecoder(resp.Body).Decode(&md)
	for _, c := range md.CubeData {
		r, err := strconv.ParseFloat(c.Rate, 64)
		if err != nil {
			return err
		}
		e.rates[c.Currency] = r
	}
	e.rates["EUR"] = 1
	// for key, element := range e.rates {
	// 	fmt.Println("Key:", key, "=>", "Element:", element)
	// }
	return nil
}

type Cubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}

type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}
