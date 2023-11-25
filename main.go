package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const DateFomat = "2006-01-02"

type RateResponse struct {
	Date         string  `json:"date"`
	BaseCurrency string  `json:"base_currency"`
	Currency     string  `json:"currency"`
	Rate         float64 `json:"rate"`
}

type Application struct {
	Data map[string]map[string]float64
}

func (a *Application) WriteError(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": msg})
}

func (a *Application) HistoricalRate(base string, currencyCode string) (float64, error) {
	var rate float64
	baseRates, exists := a.Data[base]
	if !exists {
		return rate, fmt.Errorf("base %s not found in database\n", base)
	}

	rate, exists = baseRates[currencyCode]
	if !exists {
		return rate, fmt.Errorf("Currency %s for base %s not found in database\n", currencyCode, base)
	}

	return rate, nil
}

func (a *Application) GetHistoricalRate(w http.ResponseWriter, r *http.Request) {
	fakeDelay := rand.Intn(300-200) + 200
	time.Sleep(time.Duration(fakeDelay) * time.Millisecond)
	q := r.URL.Query()
	queryDate := q.Get("date")
	if queryDate == "" {
		a.WriteError(w, "Please provide date query param (date=YYYY-MM-DD).", http.StatusBadRequest)
		return
	}

	dbDate := time.Date(2023, 9, 25, 0, 0, 0, 0, time.UTC)
	date, err := time.Parse(DateFomat, queryDate)

	if err != nil {
		a.WriteError(w, "Incorrect date value.", http.StatusBadRequest)
		return
	}

	if !date.Equal(dbDate) {
		a.WriteError(
			w,
			fmt.Sprintf("Date %s not found in database.", date.Format(DateFomat)),
			http.StatusNotFound,
		)
		return
	}

	queryBaseCurrency := q.Get("base_currency")
	if queryBaseCurrency == "" {
		a.WriteError(w, "Please provide base_currency query parameter.", http.StatusBadRequest)
		return
	}

	queryCurrency := q.Get("currency")
	if queryCurrency == "" {
		a.WriteError(w, "Please provide currency query parameter.", http.StatusBadRequest)
		return
	}

	rate, err := a.HistoricalRate(queryBaseCurrency, queryCurrency)
	if err != nil {
		a.WriteError(
			w,
			fmt.Sprintf(
				"Pair of currencies %s and %s not found in database.",
				queryBaseCurrency,
				queryCurrency,
			),
			http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&RateResponse{
		Date:         date.Format(DateFomat),
		Rate:         rate,
		BaseCurrency: queryBaseCurrency,
		Currency:     queryCurrency,
	})
}

func main() {
	app := Application{map[string]map[string]float64{
		"PLN": {
			"PLN": 1.0,
			"USD": 4.3188,
			"EUR": 4.5892,
		},
	},
	}
	port := os.Getenv("EXCHANGE_RATE_API_PORT")
	if port == "" {
		port = ":8000"
	}
	http.HandleFunc("/api/v1/historical_rates", app.GetHistoricalRate)
	fmt.Printf("Starting Exchange Rates API on port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
