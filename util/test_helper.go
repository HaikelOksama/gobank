package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"time"
)

func RandomInt(min, max int64) int64 {
	return min + rand.Int64N(max-min + 1)
}

type person struct {
	Name string `json:"name"`
}

const (
	randomuserAPI = "https://randomuser.me/api/"
)

func RandomOwner() (string, error) {	
	fmt.Printf("fetching into %v ...\n", randomuserAPI)
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		defer fmt.Printf("fetching took %v\n", duration)
	}()

	resp, err := http.Get(randomuserAPI)

	
	if err != nil {
		log.Fatal("Error: ", err)
		return "", err
	}

	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		log.Fatal("Error: ", err)
		return "", err
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal("Error: ", err)
		return "", err
	}
	
	var temp struct {
		Results []struct {
			Name struct {
				First string `json:"first"`
				Last string `json:"last"`
			} `json:"name"`
		} `json:"results"`
	}

	err = json.Unmarshal(body, &temp)
	if err != nil {
		log.Fatal("Error: ", err)
		return "", err
	}

	person := person{Name: temp.Results[0].Name.First + " " + temp.Results[0].Name.Last}

	return person.Name, nil
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "IDR", "JPY"}
	n := len(currencies)
	return currencies[rand.IntN(n)]
}