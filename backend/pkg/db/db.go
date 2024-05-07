package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/fabiokaelin/f-oauth/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var dbConn *sqlx.DB

func UpdateDBConnection() error {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.DatabaseUser, config.DatabasePassword, config.DatabaseHost, config.DatabasePort, "oauth")
	dbNew, err := sqlx.Open("mysql", connString)
	if err != nil {
		fmt.Println(err.Error())
		dbConn = nil
		return err
	}
	if dbNew == nil {
		err := errors.New("error durring updating db connection")
		dbConn = nil
		return err
	}
	// test if connection is working
	err = dbNew.Ping()
	if err != nil {
		err := errors.Join(errors.New("error durring updating db connection"), err)
		dbConn = nil
		return err
	}
	if dbConn != nil {
		dbConn.Close()
		dbConn = nil
	}
	dbConn = dbNew
	dbConn.SetMaxOpenConns(30)
	dbConn.SetMaxIdleConns(5)
	maxLifeTime := time.Minute * 30
	dbConn.SetConnMaxLifetime(maxLifeTime)
	// db.QueryRow("set client_encoding='win1252'")
	// db.QueryRow("SET CLIENT_ENCODING TO 'LATIN1';")
	return nil
}

func RunSQL(sqlStatement string, parameters ...any) (*sql.Rows, error) {
	err := dbConn.Ping()
	if err != nil {
		fmt.Println("DB Connection lost, reconnecting...")
		err := UpdateDBConnection()
		if err != nil {
			return &sql.Rows{}, err
		}
	}
	rows, err := dbConn.Query(sqlStatement, parameters...)
	if err != nil {
		err := errors.Join(errors.New("error durring executing "+sqlStatement), err)
		return &sql.Rows{}, err
	}
	return rows, nil
}

func RunSQLRow(sqlStatement string, parameters ...any) (*sql.Row, error) {

	err := dbConn.Ping()
	if err != nil {
		fmt.Println("DB Connection lost, reconnecting...")
		err := UpdateDBConnection()
		if err != nil {
			return &sql.Row{}, err
		}
	}
	rows := dbConn.QueryRow(sqlStatement, parameters...)
	return rows, nil
}

var ErrorDryExec = "DRY_EXEC is set to true, no changes will be made to the database."
