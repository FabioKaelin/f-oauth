package utils

import (
	// "backend/config"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/wpcodevo/google-github-oath2-golang/initializers"
)

var dbConn *sql.DB

func UpdateDBConnection() error {
	dbNew, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", initializers.StartConfig.DatabaseUser, initializers.StartConfig.DatabasePassword, initializers.StartConfig.DatabaseHost, initializers.StartConfig.DatabasePort, "oauth"))
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

	// spew.Dump(dbConn)
	err := dbConn.Ping()
	if err != nil {
		dbConn.Close()
		fmt.Println("1", err)
		fmt.Println("DB Connection lost, reconnecting...")
		err := UpdateDBConnection()
		if err != nil {
			fmt.Println("2", err)
			return &sql.Rows{}, err
		}
	}
	rows, err := dbConn.Query(sqlStatement, parameters...)
	if err != nil {
		newErr := errors.Join(errors.New("error durring executing "+sqlStatement), err)
		return &sql.Rows{}, newErr
	}
	return rows, nil
	// return &sql.Rows{}, errors.New(ErrorDryExec)
}

var ErrorDryExec = "DRY_EXEC is set to true, no changes will be made to the database."
