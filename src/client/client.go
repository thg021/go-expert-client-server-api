package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type DollarValue struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var r DollarValue
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Fatal(err)
	}

	SaveFile(&r)

}

func SaveFile(d *DollarValue) {
	f, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}

	tamanho, err := f.Write([]byte("Dolar: {" + d.Bid + "}"))
	if err != nil {
		panic(err)
	}
	log.Printf("File created successfully! Size: %d bytes\n", tamanho)
	f.Close()
}
