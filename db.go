package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/ini.v1"
)

type C_db struct {
	s_user      string
	s_pwd       string
	s_host      string
	s_port      string
	s_database  string
	s_engine    string
	pc_sql_db   *sql.DB
	err_row     []interface{}
	pc_sql_rows *sql.Rows
}

func (c *C_db) DB_config(cfg_file interface{}, file_section string) error {
	file, err := ini.Load(cfg_file)
	if err != nil {
		return err
	}

	section := file.Section(file_section)
	c.s_user = section.Key("user").String()
	c.s_pwd = section.Key("pwd").String()
	c.s_host = section.Key("host").String()
	c.s_port = section.Key("port").String()
	c.s_database = section.Key("database").String()
	c.s_engine = section.Key("engine").String()

	return nil
}

func (c *C_db) SQL_connection() error {
	source := c.s_user + ":" + c.s_pwd + "@tcp(" + c.s_host + ":" + c.s_port + ")/" + c.s_database
	sql_db, err := sql.Open(c.s_engine, source)
	if err != nil {
		return err
	}

	c.pc_sql_db = sql_db

	return nil
}

func (c *C_db) Insert_db(table string, row1 string, row2 string, row3 string,
	row4 string, row5 string, url string, status string, statusCode int, time string, error string) error {

	query := fmt.Sprintf("insert into %s (%s, %s, %s, %s, %s) values (?, ?, ?, ?, ?)",
		table, row1, row2, row3, row4, row5)

	_, err := c.pc_sql_db.Exec(query, url, status, statusCode, time, error)
	if err != nil {
		return err
	}

	return nil
}

func (c *C_db) Select_db(status_code_col string, table string) error {
	query := fmt.Sprintf("select * from %s where %s != ?", table, status_code_col)
	rows, err := c.pc_sql_db.Query(query, 200)
	if err != nil {
		return err
	}

	c.pc_sql_rows = rows

	for rows.Next() {
		err := rows.Scan(c.err_row)
		if err == sql.ErrNoRows {
			return nil
		}
	}

	return nil
}
