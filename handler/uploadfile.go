package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"code.google.com/p/go-uuid/uuid"

	"github.com/jeffbmartinez/log"

	"github.com/jeffbmartinez/analyzedata/storage"
)

const (
	UPLOAD_DIR = "uploads"
)

func UploadFile(response http.ResponseWriter, request *http.Request) {
	file, header, err := request.FormFile("file")
	if err != nil {
		errorMessage := "Did not receive file"

		log.Errorf("%v: %v", errorMessage, err)

		WriteSimpleResponse(response, errorMessage, http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileUuid := uuid.New()
	originalFilename := header.Filename
	originalExtension := filepath.Ext(originalFilename)
	newFilename := fmt.Sprintf("%v%v", fileUuid, originalExtension)
	storagePathname := filepath.Join(UPLOAD_DIR, newFilename)

	out, err := os.Create(storagePathname)
	if err != nil {
		log.Errorf("Unable to store the uploaded file: %v", err)

		WriteSimpleResponse(response, "", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	bytesCopied, err := io.Copy(out, file)
	if err != nil {
		log.Errorf("Couldn't copy file (%v): %v", originalFilename, err)

		WriteSimpleResponse(response, "", http.StatusInternalServerError)
		return
	}

	err = storage.PutUploadedFile(fileUuid, originalFilename, storagePathname)
	if err != nil {
		log.Errorf("Couldn't write info to db, (fileUuid: %v, originalFilename: %v, storagePathname: %v): %v",
			fileUuid, originalFilename, storagePathname, err)

		WriteSimpleResponse(response, "", http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("Uploaded %v, %v bytes", header.Filename, bytesCopied)
	log.Info(message)
	writeUploadSuccessResponse(response, fileUuid, originalFilename)
}

type uploadSuccessResponse struct {
	Uuid     string `json:"uuid"`
	Filename string `json:"filename"`
}

func writeUploadSuccessResponse(response http.ResponseWriter, uuid string, filename string) {
	message := &uploadSuccessResponse{Uuid: uuid, Filename: filename}
	WriteResponse(response, message, http.StatusOK)
}
