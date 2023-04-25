package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

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
	err_row     []interface{}
	pc_sql_db   *sql.DB
	pc_sql_rows *sql.Rows
}

func (c *C_db) Load_db_config(_s_ini_file_path string, _s_section string) error {
	file, err := ini.Load(_s_ini_file_path)
	if err != nil {
		return fmt.Errorf("Failed to read %s configuration file", _s_ini_file_path)
	}

	section, err := file.GetSection(_s_section)
	if err != nil {
		return fmt.Errorf("Failed to get %s section from %s configuration file", _s_section, _s_ini_file_path)
	}

	c.Set_db_config(section)

	return nil
}

func (c *C_db) Set_db_config(section *ini.Section) {
	c.s_user = section.Key("user").String()
	c.s_pwd = section.Key("pwd").String()
	c.s_host = section.Key("host").String()
	c.s_port = section.Key("port").String()
	c.s_database = section.Key("database").String()
	c.s_engine = section.Key("engine").String()
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

func (c *C_db) SQL_disconnect() error {
	if c.pc_sql_db != nil {
		return c.pc_sql_db.Close()
	}

	return nil
}

func (c *C_db) Create_db() error {
	create_db_query := "create database if not exists monitor"
	_, err := c.pc_sql_db.Exec(create_db_query)
	if err != nil {
		return err
	}

	var table string
	row := c.pc_sql_db.QueryRow(`select table_name from information_schema.tables 
	where table_name like "http_server%" order by table_name desc limit 1`)
	if err := row.Scan(&table); err != nil {
		if err == sql.ErrNoRows {
			table = "http_server_1"
		} else {
			return err
		}
	} else {
		last_num, err := strconv.Atoi(strings.TrimPrefix(table, "http_server_"))
		if err != nil {
			return err
		}
		table = fmt.Sprintf("http_server_%d", last_num+1)
	}

	create_table_query := fmt.Sprintf(`create table if not exists %s (
		id int not null auto_increment primary key
	);`, table)
	_, err = c.pc_sql_db.Exec(create_table_query)
	if err != nil {
		return err
	}

	add_col_query := fmt.Sprintf(`alter table %s add (
		url varchar(255) not null,
		status char(50) not null,
		status_code int(10) not null,
		time varchar(50) not null,
		error varchar(255) default null
	)`, table)
	_, err = c.pc_sql_db.Exec(add_col_query)
	if err != nil {
		return err
	}

	return nil
}

func (c *C_db) Insert_data(_s_table string, _s_url string, _s_status string, _i_status_code int, _s_time string, _s_err string) error {
	add_row_query := fmt.Sprintf(`insert into %s (url, status, status_code, time, error) 
	values (?, ?, ?, ?, ?)`, _s_table)

	_, err := c.pc_sql_db.Exec(add_row_query, _s_url, _s_status, _i_status_code, _s_time, _s_err)
	if err != nil {
		return err
	}

	return nil
}

func (c *C_db) Select_data() error {
	tables := []string{}
	rows, err := c.pc_sql_db.Query(`select table_name from information_schema.tables 
	where table_schema = "monitor"`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		if err != nil {
			return err
		}
		tables = append(tables, table)
	}

	for _, table := range tables {
		query := fmt.Sprintf("select * from %s where status_code != ?", table)
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
	}

	return nil
}

func (c *C_db) Close_rows() error {
	if c.pc_sql_rows != nil {
		return c.pc_sql_rows.Close()
	}

	return nil
}
