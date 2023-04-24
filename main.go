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

	time_interval := time.Second * 5 // 모니터링 시간 간격
	ticker := time.NewTicker(time_interval)
	defer ticker.Stop()

	for range ticker.C {

		urls := []string{"http://naver.com", "http://google.com"} //대상 url 입력
		h_list := []C_http{}

		for _, url := range urls {
			err := h.Get_http_status(url)
			if err != nil {
				log.Printf("errer: %v\n", err)
				os.Exit(1)
			}
			h_list = append(h_list, h)
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

		for _, list := range h_list {
			err = d.Insert_db("http_server", "URL", "status", "status_code", "time", "error",
				list.s_url, list.s_status, list.i_status_code, list.s_time, list.s_error)
			if err != nil {
				log.Printf("errer: %v\n", err)
				os.Exit(1)
			}
		}

		err = d.Select_db("status_code", "http_server")
		if err != nil {
			log.Printf("errer: %v\n", err)
			os.Exit(1)
		}
		defer d.pc_sql_rows.Close()

		for _, v := range d.err_row {
			fmt.Printf("%d \n", v)
		}

		if len(d.err_row) == 0 {
			s.Init("ap-northeast-2", "AKIAVOZYFWFTBWEOBG7T",
				"eHqPu4vSNraNS9IYNF7dnuKPI7vSSR8OXFuvzPyN")
			// (email) 보내는 사람, 받는 사람, 제목, 내용 입력
			sender := "abh4017@naver.com"
			recipients := []string{
				"abh4017@naver.com",
				"qudgusyou012@gmail.com",
			}
			subject := "monitoring error"
			body := ""
			for _, v := range d.err_row {
				body += fmt.Sprintf("error rows: %s\n", v.(string))
			}

			s.Write_email(sender, recipients, subject, body)

			err := s.Set_config()
			if err != nil {
				log.Printf("errer: %v\n", err)
				os.Exit(1)
			}

			err = s.Send_email(s.pc_client, s.s_sender, s.s_recipient, s.s_subject, s.s_body)
		}
	}
}
