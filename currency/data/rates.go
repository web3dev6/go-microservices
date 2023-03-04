package data

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/go-hclog"
)

type ExchangeRates struct {
	log   hclog.Logger
	rates map[string]float64
}

// MonitorRates checks the rates in the ECB API every interval and sends a message to the
// returned channel when there are changes
//
// Note: the ECB API only returns data once a day, this function only simulates the changes
// in rates for demonstration purposes
func (e *ExchangeRates) MonitorRates(interval time.Duration) chan struct{} {
	ret := make(chan struct{})
	go func() {
		ticker := time.NewTicker(interval)
		for {
			<-ticker.C
			// just add a random difference to the rate and return it
			// this simulates the fluctuations in currency rates
			for k, v := range e.rates {
				// change can be 10% of original value
				change := (rand.Float64() / 10)
				// is this a postive or negative change
				direction := rand.Intn(2)
				if direction == 0 {
					// new value with be min 90% of old
					change = 1 - change
				} else {
					// new value will be 110% of old
					change = 1 + change
				}
				// modify the rate
				e.rates[k] = v * change
			}
			// notify updates, this will block unless there is a listener on the other end
			ret <- struct{}{}
		}
	}()
	return ret
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
	err := er.fetchRatesFromECB()
	return er, err
}

// GetRates fetches currency rates against EUR for various currencies from eu-central-bank
func (e *ExchangeRates) fetchRatesFromECB() error {
	resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		return err
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
