package main

import (
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

func (c *C_http) Get_http_status(url string) (status_code int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	c.s_url = url
	c.i_status_code = resp.StatusCode
	c.s_status = resp.Status
	c.s_time = time.Now().Format("2006-01-02 15:04:05")

	return c.i_status_code, nil
}
