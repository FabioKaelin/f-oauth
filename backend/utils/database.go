package utils

import (
	// "backend/config"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/wpcodevo/google-github-oath2-golang/initializers"
)

func getDBconnection() *sql.DB {
	// fmt.Println("connect to Database")
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", initializers.StartConfig.DatabaseUser, initializers.StartConfig.DatabasePassword, initializers.StartConfig.DatabaseHost, initializers.StartConfig.DatabasePort, "oauth"))
	// db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", "root", "secretpass", "db.fabkli.ch", "38487", "oauth"))
	// db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DB_Username, DB_Password, DB_Host, DB_Port, DB_Name))
	if err != nil {
		fmt.Println(err.Error())
	}
	// db.QueryRow("set client_encoding='win1252'")
	// db.QueryRow("SET CLIENT_ENCODING TO 'LATIN1';")
	return db
}

func RunSQL_OLD(sqlStatement string) *sql.Rows {
	db := getDBconnection()
	fmt.Println("run sql")
	rows, err := db.Query(sqlStatement)
	defer db.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	// defer rows.Close()
	return rows
}

func RunSQLRow(sqlStatement string) *sql.Row {
	db := getDBconnection()
	fmt.Println("run sql row")
	row := db.QueryRow(sqlStatement)
	defer db.Close()
	// defer rows.Close()
	return row
}

func RunSQLSecureOne(sqlStatement string, parameters ...any) (*sql.Rows, error) {
	db := getDBconnection()
	fmt.Println("run sql")
	rows, err := db.Query(sqlStatement, parameters...)
	defer db.Close()
	if err != nil {
		return new(sql.Rows), err
	}
	// defer rows.Close()
	return rows, nil
}

func RunSQLSecureMultible(sqlstatements [][]string) {
	fmt.Println("run sql")
	db := getDBconnection()
	for _, statement := range sqlstatements {
		var parametersAny []any
		for _, parameter := range statement[1:] {
			parametersAny = append(parametersAny, parameter)
		}
		rows, err := db.Query(statement[0], parametersAny...)
		defer db.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
		rows.Close()
	}
}
