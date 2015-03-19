package handler

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jeffbmartinez/log"

	"github.com/jeffbmartinez/analyzedata/storage"
)

func DownloadFile(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	uuid := vars["uuid"]

	log.Infof("Requested download of file with uuid='%v'", uuid)

	uploadedFile, err := storage.GetUploadedFile(uuid)
	switch {
	case err == storage.UploadedFileNotFound:
		log.Infof("File with uuid='%v' not found, doesn't exist", uuid)
		WriteSimpleResponse(response, "Requested file not found", http.StatusNotFound)
		return
	case err != nil:
		log.Errorf("File with uuid='%v' not found, something went wrong", uuid)
		WriteSimpleResponse(response, "", http.StatusInternalServerError)
		return
	}

	log.Infof("File with uuid='%v' found", uuid)

	http.ServeFile(response, request, uploadedFile.StoragePath)
}
