package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func fileReader() string {
	content, err := ioutil.ReadFile("C:\\Work\\GO\\RSO_projekt\\mapsapi\\KEYS.TXT")

	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}
func main() {
	APIKEY := fileReader()
	fmt.Println(APIKEY)
	origin := "Ptuj"
	destination := "Ljubljana|Maribor|portorož"
	url := "https://maps.googleapis.com/maps/api/distancematrix/json?origins=" + origin +
		"&destinations=" + destination + "&units=metricsapi&key=" + APIKEY + "&avoidHighways=True"

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
