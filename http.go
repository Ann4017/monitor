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
	pc_response   *http.Response
}

func (c *C_http) Get_http_status(_s_url string) error {
	client := &http.Client{
		Timeout: time.Duration(time.Second * 15),
	}

	resp, err := client.Get(_s_url)
	if err != nil {
		c.s_error = err.Error()
		return err
	}

	c.pc_response = resp

	c.s_url = _s_url
	c.i_status_code = resp.StatusCode
	c.s_status = resp.Status
	c.s_time = time.Now().Format("2006-01-02 15:04:05")

	fmt.Println(c.s_url)
	fmt.Println(c.i_status_code)
	fmt.Println(c.s_status)
	fmt.Println(c.s_time)

	return nil
}

func (c *C_http) Close_resp_body() error {
	if c.pc_response != nil {
		return c.pc_response.Body.Close()
	}

	return nil
}
