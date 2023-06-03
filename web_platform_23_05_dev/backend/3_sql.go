package backend

//

// IMPORTS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"database/sql"
	"fmt"

	// _ "github.com/lib/pq"
	_ "github.com/sijms/go-ora/v2"
)

// IMPORTS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

//

// STRUCTS ////////////////////////////////////////////////////////////////////////////////////////////////////////////
// STRUCTS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

//

// METHODS ////////////////////////////////////////////////////////////////////////////////////////////////////////////
// METHODS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

//

// GLOBALS ////////////////////////////////////////////////////////////////////////////////////////////////////////////
// GLOBALS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

//

// ACTIONS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

func CreateDbPgConnection() (*sql.DB, error) {
	source := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		ConfigV.HostPgSQL, ConfigV.PortPgSQL, ConfigV.UserPgSQL,
		ConfigV.PasswordPgSQL, ConfigV.DatabasePgSQL, ConfigV.SslModePgSQL,
	)
	dbConnection, err := sql.Open(ConfigV.DriverNamePgSQL, source)
	if err != nil {
		return nil, err
	}

	return dbConnection, nil
}

func ExecuteSelectOneDb(object []any, query string, args ...any) error {
	dbConnection, err := CreateDbPgConnection()
	if err != nil {
		return err
	}
	defer func(dbConnection *sql.DB) {
		err = dbConnection.Close()
		if err != nil {
			return
		}
	}(dbConnection)

	rows, err := dbConnection.Query(query, args...)
	if err != nil {
		return err
	}

	rows.Next()
	err = rows.Scan(object...)
	if err != nil {
		return err
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

func ExecuteSelectManyDb(objects *[]string, query string, args ...any) error {
	dbConnection, err := CreateDbPgConnection()
	if err != nil {
		return err
	}
	defer func(dbConnection *sql.DB) {
		err = dbConnection.Close()
		if err != nil {
			return
		}
	}(dbConnection)

	rows, err := dbConnection.Query(query, args...)
	if err != nil {
		return err
	}

	usernames := make([]string, 0)
	for rows.Next() {
		var obj string
		err = rows.Scan(&obj)
		if err != nil {
			return err
		}

		usernames = append(usernames, obj)
	}
	fmt.Println(usernames)
	*objects = usernames

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

func ExecuteRowsDb(query string, args ...any) (*sql.Rows, error) {
	dbConnection, err := CreateDbPgConnection()
	if err != nil {
		return nil, err
	}
	defer func(dbConnection *sql.DB) {
		err = dbConnection.Close()
		if err != nil {
			return
		}
	}(dbConnection)

	rows, err := dbConnection.Query(query, args...)
	if err != nil {
		return nil, err
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func ExecuteInsertOrDeleteDb(query string, args ...any) error {
	dbConnection, err := CreateDbPgConnection()
	if err != nil {
		return err
	}
	defer func(dbConnection *sql.DB) {
		err = dbConnection.Close()
		if err != nil {
			return
		}
	}(dbConnection)

	dbTransaction, err := dbConnection.Begin()
	if err != nil {
		return err
	}
	defer func(dbTransaction *sql.Tx) {
		_ = dbTransaction.Rollback()
	}(dbTransaction)

	_, err = dbTransaction.Exec(query, args...)
	if err != nil {
		return err
	}

	err = dbTransaction.Commit()
	if err != nil {
		return err
	}

	return nil
}

// ACTIONS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

//

// EXTRAS /////////////////////////////////////////////////////////////////////////////////////////////////////////////
// EXTRAS /////////////////////////////////////////////////////////////////////////////////////////////////////////////

//
