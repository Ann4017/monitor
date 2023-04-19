package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/ini.v1"
)

type C_db struct {
	s_user     string
	s_pwd      string
	s_host     string
	s_port     string
	s_database string
	s_engine   string
	pc_sql_db  *sql.DB
	err_row    []interface{}
}

var h = C_http{}

func (c *C_db) DB_config(cfg_file interface{}, file_section string) error {
	file, err := ini.Load(cfg_file)
	if err != nil {
		return err
	}

	section := file.Section(file_section)
	c.s_user = section.Key("root").String()
	c.s_pwd = section.Key("RLGH3qjs!!").String()
	c.s_host = section.Key("localhost").String()
	c.s_port = section.Key("3306").String()
	c.s_database = section.Key("http").String()
	c.s_engine = section.Key("mysql").String()

	return nil
}

func (c *C_db) SQL_connection() error {
	source := c.s_user + ":" + c.s_pwd + "@tcp(" + c.s_host + ":" + c.s_port + ")/" + c.s_database
	sql_db, err := sql.Open(c.s_engine, source)
	if err != nil {
		return err
	}

	c.pc_sql_db = sql_db
	defer sql_db.Close()

	return nil
}

func (c *C_db) Insert_db(table string, row1 string, row2 string, row3 string,
	row4 string, row5 string) error {

	str1 := fmt.Sprintf("insert into %s (%s, %s, %s, %s, %s) ",
		table, row1, row2, row3, row4, row5)
	str2 := fmt.Sprintf("values (%s, %s, %d, %s, %s)",
		h.s_url, h.s_status, h.i_status_code, h.s_time, h.s_error)
	query := str1 + str2

	_, err := c.pc_sql_db.Exec(query)
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
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(c.err_row)
		if err == sql.ErrNoRows {
			return nil
		}
	}

	return nil
}
