package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	h := C_http{}
	d := C_db{}
	s := C_ses{}

	log_file, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
	}

	defer log_file.Close()

	log.SetOutput(log_file)

	for {
		_, err := h.Get_http_status("https://www.google.com/")
		if err != nil {
			log.Printf("errer: %v\n", err)
			os.Exit(1)
		}

		err = d.DB_config("config.ini", "database")
		if err != nil {
			log.Printf("errer: %v\n", err)
			os.Exit(1)
		}
		defer d.pc_sql_db.Close()

		err = d.SQL_connection()
		if err != nil {
			log.Printf("errer: %v\n", err)
			os.Exit(1)
		}

		err = d.Insert_db("http_server", "URL", "status", "status_code", "time", "error")
		// h.s_url, h.s_status, h.i_status_code, h.s_time, h.s_error)
		if err != nil {
			log.Printf("errer: %v\n", err)
			os.Exit(1)
		}

		err = d.Select_db("status_code", "http_server")
		if err != nil {
			log.Printf("errer: %v\n", err)
			os.Exit(1)
		}
		defer d.pc_sql_rows.Close()

		if len(d.err_row) != 0 {
			s.Init("ap-northeast-2", "AKIAVOZYFWFTBWEOBG7T",
				"eHqPu4vSNraNS9IYNF7dnuKPI7vSSR8OXFuvzPyN")

			s.Write_email("abh4017@naver.com", "abh4017@naver.com", "monitor test", "test")

			err := s.Set_config()
			if err != nil {
				log.Printf("errer: %v\n", err)
				os.Exit(1)
			}

			err = s.Send_email(s.pc_client, s.s_sender, s.s_recipient, s.s_subject, s.s_body)
		}

		time.Sleep(time.Second * 10)
	}
}
