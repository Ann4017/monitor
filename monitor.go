package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Monitor struct {
	h C_http
	d C_db
	s C_ses
}

func (m *Monitor) Create_log() error {
	log_file, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer log_file.Close()

	log.SetOutput(log_file)

	return nil
}

func (m *Monitor) Init(_s_ini_file_path, _s_section, _s_region, _s_access_key, _s_secret_key string) error {
	err := m.Create_log()
	if err != nil {
		return err
	}

	err = m.d.Load_db_config(_s_ini_file_path, _s_section)
	if err != nil {
		return err
	}

	err = m.d.SQL_connection()
	if err != nil {
		return err
	}

	m.s.Init(_s_region, _s_access_key, _s_secret_key)

	err = m.s.Set_config()
	if err != nil {
		return err
	}

	return nil
}

func (m *Monitor) Run(_s_urls []string, _s_sender string, _s_recipient []string,
	monitoring_second_cycle time.Duration) error {
	table, err := m.d.Create_db()
	if err != nil {
		return err
	}

	time_interval := time.Second * monitoring_second_cycle
	ticker := time.NewTicker(time_interval)
	defer ticker.Stop()

	for range ticker.C {
		urls_info := []C_http{}

		for _, url := range _s_urls {
			err := m.h.Get_http_status(url)
			if err != nil {
				return err
			}
			urls_info = append(urls_info, m.h)
		}

		for _, info := range urls_info {
			err := m.d.Insert_data(table, info.s_url, info.s_status, info.i_status_code, info.s_time, info.s_error)
			if err != nil {
				return err
			}
		}

		err = m.d.Select_data()
		if err != nil {
			return err
		}

		if len(m.d.err_row) != 0 {
			sender := _s_sender
			recipients := _s_recipient
			subject := "Finding problems while monitoring"
			body := ""

			for _, v := range m.d.err_row {
				body += fmt.Sprintf("error rows: %s\n", v.(string))
			}

			m.s.Write_email(sender, recipients, subject, body)

			err := m.s.Send_email(m.s.pc_client, sender, recipients, subject, body)
			if err != nil {
				return err
			}
		}
	}

	defer m.h.Close_resp_body()
	defer m.d.SQL_disconnect()
	defer m.d.Close_rows()

	return nil
}
