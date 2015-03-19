// Copyright 2015 Jeff Martinez. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE.txt file
// or at http://opensource.org/licenses/MIT

package storage

import (
	"database/sql"
	"errors"

	"github.com/jeffbmartinez/log"
)

func PutUploadedFile(uuid string, originalFilename string, storagePathname string) error {
	log.Infof("About to put uploaded file info (%v, %v, %v) in db",
		uuid, originalFilename, storagePathname)

	db, err := GetDbConnection()
	if err != nil {
		log.Error("Unable to store the uploaded file info in the db, can't access db")
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO uploads (uuid, original_filename, storage_path) VALUES (?, ?, ?)",
		uuid, originalFilename, storagePathname)
	if err != nil {
		log.Errorf("Problem storing upload info to db (uuid: %v, original_filename: %v, storage_path: %v): %v",
			uuid, originalFilename, storagePathname, err)
	}

	log.Infof("Uploaded file info (uuid: %v) stored", uuid)

	return err
}

func GetUploadedFile(uuid string) (Upload, error) {
	upload := Upload{Uuid: uuid}

	log.Infof("About to retrieve uploaded file info (uuid: %v) from db", upload.Uuid)

	db, err := GetDbConnection()
	if err != nil {
		log.Error("Unable to retrieve the file info from the db, can't access db")
		return upload, err
	}
	defer db.Close()

	query :=
		`SELECT id, original_filename, storage_path, upload_date
		 FROM uploads
		 WHERE uuid=?`

	row := db.QueryRow(query, upload.Uuid)

	err = row.Scan(&upload.Id, &upload.OriginalFilename, &upload.StoragePath, &upload.UploadDate)
	switch {
	case err == sql.ErrNoRows:
		log.Infof("Found no user with uuid='%v'", uuid)
		return upload, UploadedFileNotFound
	case err != nil:
		log.Errorf("Unknown error while retrieving upload file info with uuid='%v': %v", upload.Uuid, err)
		return upload, err
	}

	log.Info("Retrieved file info (uuid: %v)", upload.Uuid)

	return upload, err
}

type Upload struct {
	Id               int
	Uuid             string
	OriginalFilename string
	StoragePath      string
	UploadDate       string
}

var UploadedFileNotFound = errors.New("Could not find the requested file")
