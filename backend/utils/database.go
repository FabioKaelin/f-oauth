package utils

import (
	// "backend/config"
	"database/sql"
	"errors"
	"fmt"

	"github.com/fabiokaelin/f-oauth/initializers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var dbConn *sqlx.DB

func UpdateDBConnection() error {
	dbNew, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", initializers.StartConfig.DatabaseUser, initializers.StartConfig.DatabasePassword, initializers.StartConfig.DatabaseHost, initializers.StartConfig.DatabasePort, "oauth"))
	if err != nil {
		fmt.Println(err.Error())
	}
	// test if connection is working
	err = dbNew.Ping()
	if err != nil {
		newErr := errors.Join(errors.New("error durring updating db connection"), err)
		return newErr
	}
	dbConn = dbNew
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
		newErr := errors.Join(errors.New("error durring executing "+sqlStatement), err)
		return &sql.Rows{}, newErr
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
