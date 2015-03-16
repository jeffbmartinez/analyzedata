package handler

import (
	"encoding/json"
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

		writeSimpleResponse(response, errorMessage, http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileUuid := uuid.New()
	originalFilename := header.Filename
	originalExtension := filepath.Ext(originalFilename)
	storagePathname := fmt.Sprintf("%v/%v%v", UPLOAD_DIR, fileUuid, originalExtension)

	out, err := os.Create(storagePathname)
	if err != nil {
		log.Errorf("Unable to store the uploaded file: %v", err)

		writeSimpleResponse(response, "", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	bytesCopied, err := io.Copy(out, file)
	if err != nil {
		log.Errorf("Couldn't copy file (%v): %v", originalFilename, err)

		writeSimpleResponse(response, "", http.StatusInternalServerError)
		return
	}

	err = storage.StoreUpload(fileUuid, originalFilename, storagePathname)
	if err != nil {
		log.Errorf("Couldn't write info to db, (fileUuid: %v, originalFilename: %v, storagePathname: %v): %v",
			fileUuid, originalFilename, storagePathname, err)

		writeSimpleResponse(response, "", http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("Uploaded %v, %v bytes", header.Filename, bytesCopied)
	log.Info(message)
	writeUploadSuccessResponse(response, fileUuid, originalFilename)
}

type simpleResponse struct {
	Message string `json:"message"`
}

type uploadSuccessResponse struct {
	Uuid     string `json:"uuid"`
	Filename string `json:"filename"`
}

func writeSimpleResponse(response http.ResponseWriter, message string, statusCode int) {
	msg := &simpleResponse{Message: message}
	writeResponse(response, msg, statusCode)
}

func writeUploadSuccessResponse(response http.ResponseWriter, uuid string, filename string) {
	message := &uploadSuccessResponse{Uuid: uuid, Filename: filename}
	writeResponse(response, message, http.StatusOK)
}

func writeResponse(response http.ResponseWriter, message interface{}, statusCode int) {
	responseString, err := json.Marshal(message)

	if err != nil {
		log.Errorf("Couldn't marshal json: %v", err)

		response.WriteHeader(statusCode)
		response.Write([]byte(""))
		return
	}

	response.WriteHeader(statusCode)
	response.Write([]byte(responseString))
}
