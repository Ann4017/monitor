package main

import (
	"log"
	"testing"
)

func Test(t *testing.T) {
	m := Monitor{}
	urls := []string{"http://naver.com", "http://google.com"}
	recipients := []string{"abh4017@naver.com", "qudgusyou012@gmail.com"}

	err := m.Init("config.ini", "database", "", "", "")
	if err != nil {
		log.Fatal(err)
	}

	err = m.Run(urls, "abh4017@naver.com", recipients, 10)
	if err != nil {
		log.Fatal(err)
	}
}
