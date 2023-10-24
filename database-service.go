package databaseservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	//_ "github.com/denisenkom/go-mssqldb"
)

type DatabaseService struct {
	Database *gorm.DB
	//Database *sql.DB
}

func Init() *DatabaseService {
	dbHost := "LAPTOP-7J78SQ4H"
	dbUser := "krisna"
	dbPass := "123"
	dbName := "transaction"
	dbPort := "1433"

	connectionStr := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;",
		dbHost, dbUser, dbPass, dbPort, dbName)

	// db, err := sql.Open("mssql", connectionStr)
	// if err != nil {
	// 	fmt.Println("Error connecting to database:", err)
	// }

	db, err := gorm.Open(sqlserver.Open(connectionStr), &gorm.Config{})
	if err != nil {
		fmt.Println("Error connecting to database:", err)
	}

	return &DatabaseService{Database: db}

}

func (ds *DatabaseService) CloseDBConnection() error {
	if ds.Database == nil {
		return errors.New("Database connection is nil")
	}

	db, err := ds.Database.DB()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
	}

	err = db.Close()
	if err != nil {
		return err
	}

	return nil
}

// func (ds *DatabaseService) CallStoredProcedure(c context.Context, spName string, params map[string]any) (*sql.Rows, error) {
// 	placeholders := make([]string, len(params))
// 	paramValues := make([]interface{}, len(params))

// 	i := 0
// 	for key, val := range params {
// 		placeholders[i] = fmt.Sprintf("%s = ?", key)
// 		paramValues[i] = val
// 		i++
// 	}

// 	query := spName + " " + strings.Join(placeholders, ",")

// 	stmt, err := ds.Database.PrepareContext(c, query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer stmt.Close()

// 	result, err := stmt.QueryContext(c, paramValues...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }

func (ds *DatabaseService) CallStoredProcedure(c context.Context, spName string, params map[string]any) (*sql.Rows, error) {
	placeholders := make([]string, len(params))
	paramValues := make([]interface{}, len(params))

	i := 0
	for key, val := range params {
		placeholders[i] = fmt.Sprintf("%s = @param%d", key, i+1)
		paramValues = append(paramValues, sql.Named(fmt.Sprintf("param%d", i+1), val))
		i++
	}

	query := fmt.Sprintf("EXEC %s %s", spName, strings.Join(placeholders, ","))

	db, err := ds.Database.DB()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
	}

	stmt, err := db.PrepareContext(c, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(c, paramValues...)
	if err != nil {
		return nil, err
	}

	return result, nil
}
