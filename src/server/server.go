package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ExchangeRate struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

type DollarValue struct {
	ID  uuid.UUID `gorm:"type:uuid;primaryKey;" json:"-"`
	Bid string    `json:"bid"`
}

const DB_MAX_TIMEOUT = 10 * time.Millisecond
const REQUEST_MAX_DURATION = 200 * time.Millisecond

func main() {
	http.HandleFunc("/cotacao", ExchangeRatesHandler)
	http.ListenAndServe(":8080", nil)
}

func ExchangeRatesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting search exchange rates")
	exchangeRateResponse, err := searchExchangeRate()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	v, err := saveExchangeRate(exchangeRateResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(v)
}

func searchExchangeRate() (*ExchangeRate, error) {

	ctxHttp, cancel := context.WithTimeout(context.Background(), REQUEST_MAX_DURATION)
	defer cancel()
	log.Println("Starting get request")
	req, err := http.NewRequestWithContext(ctxHttp, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		log.Panic("Error in the request to the dollar quotation service")
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var r ExchangeRate
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil

}

func saveExchangeRate(exchangeRateResponse *ExchangeRate) (*DollarValue, error) {
	log.Printf("Opening connection with db")

	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Panic("Opening database connection")
	}
	db.AutoMigrate(&DollarValue{})
	gormCtx, gormCancel := context.WithTimeout(context.Background(), DB_MAX_TIMEOUT)
	defer gormCancel()

	dollarValue := &DollarValue{
		ID:  uuid.NewV4(),
		Bid: exchangeRateResponse.Usdbrl.Bid,
	}
	result := db.WithContext(gormCtx).Create(dollarValue)

	if result.Error != nil {
		log.Panic("Error creating exchange rate")
	}
	log.Printf("Data saved to database successfully")

	return dollarValue, nil
}
