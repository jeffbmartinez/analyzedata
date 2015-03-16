// Copyright 2015 Jeff Martinez. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE.txt file
// or at http://opensource.org/licenses/MIT

package storage

import (
	"fmt"

	"github.com/jeffbmartinez/log"
)

func StoreUpload(uuid string, originalFilename string, storagePathname string) error {
	db, err := GetDbConnection()
	if err != nil {
		log.Error("Unable to store the uploaded file info in the db, can't access db")
		return err
	}
	defer db.Close()

	fmt.Printf("Would be storing (%v, %v, %v) to db\n", uuid, originalFilename, storagePathname)

	return nil
}
