package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"io/ioutil"
	"net/http"
)

type jsn struct {
	// body struct
	Result struct {
		Bencin string `json:"gasoline"`
		Dizel  string `json:"diesel"`
	} `json:"result"`
}
type jsnret struct {
	Bencin string `json:"bencin"`
	Dizel  string `json:"dizel"`
}

func getDataJson() {
	url := "https://api.collectapi.com/gasPrice/fromCity?city=ljubljana?currency=eur'"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", "apikey 2M0y9SHvCFNV5KUD2lGZL2:3VnJ9JIwyF4UCf01Ffbx3S")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("empty")
		fmt.Println(err)

	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)

	}
	var data jsn
	var retrn jsnret
	if err := json.Unmarshal(body, &data); err != nil { // Parse []byte to go struct pointer
		fmt.Println(err)
		fmt.Println("Can not unmarshal JSON")
	}
	retrn.Dizel = data.Result.Dizel
	retrn.Bencin = data.Result.Bencin
	fmt.Println(retrn)
	vrni, err := json.Marshal(retrn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(vrni))
	//return vrni from function

}
func main() {
	//getDataJson()
	app := fiber.New()

	app.Use(cors.New())

	app.Get("/api/gas", func(c *fiber.Ctx) error {
		url := "https://api.collectapi.com/gasPrice/fromCity?city=ljubljana?currency=eur'"

		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Add("content-type", "application/json")
		req.Header.Add("authorization", "apikey 2M0y9SHvCFNV5KUD2lGZL2:3VnJ9JIwyF4UCf01Ffbx3S")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("empty")
			fmt.Println(err)

		}

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)

		}
		var data jsn
		var retrn jsnret
		if err := json.Unmarshal(body, &data); err != nil { // Parse []byte to go struct pointer
			fmt.Println(err)
			fmt.Println("Can not unmarshal JSON")
		}
		retrn.Dizel = data.Result.Dizel
		retrn.Bencin = data.Result.Bencin
		fmt.Println(retrn)
		vrni, err := json.Marshal(retrn)
		if err != nil {
			fmt.Println(err)
		}
		return c.Send(vrni)
	})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Send([]byte("Gasoline api container working"))
	})

	app.Listen(":8004")

}
