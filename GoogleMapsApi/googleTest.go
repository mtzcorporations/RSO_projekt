package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	origin := "Ptuj"
	destination := "Ljubljana|Maribor|portorož  "
	APIKEY := "AIzaSyDc_d1tvLqZhDC3sBhgtBh5DxMVJGMajps"
	url := "https://maps.googleapis.com/maps/api/distancematrix/json?origins=" + origin +
		"&destinations=" + destination + "&units=metrics&key=" + APIKEY + "&avoidHighways=True"
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
