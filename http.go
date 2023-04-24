package main

import (
	"fmt"
	"net/http"
	"time"
)

type C_http struct {
	s_url         string
	s_status      string
	i_status_code int
	s_time        string
	s_error       string
}

func (c *C_http) Get_http_status(_s_url string) (err error) {
	client := &http.Client{
		Timeout: time.Duration(time.Second * 10),
	}

	resp, err := client.Get(_s_url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	c.s_url = _s_url
	c.i_status_code = resp.StatusCode
	c.s_status = resp.Status
	c.s_time = time.Now().Format("2006-01-02 15:04:05")

	fmt.Println(c.s_url)
	fmt.Println(c.i_status_code)
	fmt.Println(c.s_status)
	fmt.Println(c.s_time)
	fmt.Println(c.s_error)

	return nil
}
