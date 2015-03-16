// Copyright 2015 Jeff Martinez. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE.txt file
// or at http://opensource.org/licenses/MIT

package storage

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // Register 'sqlite3' driver

	"github.com/jeffbmartinez/log"
)

const (
	SQLITE3_DRIVER = "sqlite3"
	DB_LOCATION    = "data/analyzedata.sqlite"
)

/*
Don't forget to close the connection when you're done with it
	defer db.Close()
*/
func GetDbConnection() (*sql.DB, error) {
	db, err := sql.Open(SQLITE3_DRIVER, DB_LOCATION)
	if err != nil {
		log.Errorf("Can't access db (driver: %v, db_location: %v): %v",
			SQLITE3_DRIVER, DB_LOCATION, err)
	}

	return db, err
}
