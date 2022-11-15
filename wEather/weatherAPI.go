package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type weather struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
	Oblaki []struct {
		Sonce string `json:"main"`
	} `json:"weather"`
}

func main() {
	url := "http://api.openweathermap.org/data/2.5/weather?lat=46.05&lon=14.50&units=metric&appid=ab8428d16bce2694fb18fbab32071873"
	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "test")
	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	vreme := weather{}
	jsonErr := json.Unmarshal([]byte(body), &vreme)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	fmt.Println(vreme.Name, vreme.Main, vreme.Oblaki[0])

}
