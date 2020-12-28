package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/evrenios/letmein/misc"
)

func main() {
	var hour int
	var endpoint string
	flag.IntVar(&hour, "hour", 1, "")
	flag.StringVar(&endpoint, "endpoint", "", "Specifcy the endpoint that your service is listening")
	flag.Parse()

	if len(endpoint) == 0 {
		fmt.Println("endpoint can not be blank")
		os.Exit(0)
	}

	if hour < 0 || hour > 23 {
		fmt.Println("hour must be between 1 and 23")
		os.Exit(0)
	}

	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	resp, err := http.Get("http://whatismyip.akamai.com")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bIP, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	req := &misc.AuthReq{
		IP:     string(bIP),
		Secret: misc.Secret,
		Hour:   hour,
		Name:   name,
	}
	b, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	if _, err = http.Post(endpoint, "application/json", bytes.NewReader(b)); err != nil {
		panic(err)
	}
	fmt.Println("your ip has been authorized for an hour")
}
