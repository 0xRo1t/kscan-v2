package mssql

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"time"

	mssql "github.com/microsoft/go-mssqldb"
)

func Check(Host, Username, Password string, Port int) error {
	mssql.SetLogger(log.New(io.Discard, "", log.Ldate|log.Ltime))
	dataSourceName := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%v;encrypt=disable;timeout=%v", Host, Username, Password, Port, 5*time.Second)
	db, err := sql.Open("mssql", dataSourceName)
	if err != nil {
		return err
	}
	db.SetConnMaxLifetime(5 * time.Second)
	db.SetMaxIdleConns(0)
	defer db.Close()
	err = db.Ping()
	if err != nil {
		return err
	}
	return nil
}
