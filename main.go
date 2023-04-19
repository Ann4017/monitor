package main

import (
	"fmt"
	"net/http"
	"time"
)

type HTTP struct {
	status int
}

func (c *HTTP) Get_http_status(url string) (status_code int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	c.status = resp.StatusCode
	time.Sleep(time.Second)

	return c.status, nil
}

func main() {
	c := HTTP{}

	for {
		stutus, err := c.Get_http_status("https://www.google343212.co.kr/")
		if err != nil {
			if stutus == 0 {
				fmt.Println("No such host")
			} else if stutus != 200 {
				fmt.Printf("Error status code: %d\n", stutus)
			}
		}
		fmt.Printf("Current Status code: %d\n", stutus)
		time.Sleep(time.Second)
	}
}
