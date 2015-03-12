package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type responseMessage struct {
	Message string `json:"message,omitempty"`
}

func UploadFile(response http.ResponseWriter, request *http.Request) {
	msg := &responseMessage{}

	file, header, err := request.FormFile("file")

	if err != nil {
		log.Println("[ERROR] Couldn't grab uploaded file info", err)

		response.WriteHeader(http.StatusBadRequest)
		msg.Message = fmt.Sprint(err)
	} else {
		msg.Message = fmt.Sprintf("Uploaded %v successfully", header.Filename)
	}

	defer file.Close()

	out, err := os.Create("/tmp/uploadedfile")
	if err != nil {
		errorMessage := "Unable to store the uploaded file"

		log.Println("[ERROR]", errorMessage, err)

		response.WriteHeader(http.StatusInternalServerError)
		msg.Message = fmt.Sprint(errorMessage)
	}

	defer out.Close()

	bytesCopied, err := io.Copy(out, file)
	if err != nil {
		errorMessage := "Couldn't copy file"

		log.Println("[ERROR]", errorMessage, err)

		response.WriteHeader(http.StatusInternalServerError)
		msg.Message = fmt.Sprintf(errorMessage)
	} else {
		log.Printf("[INFO] Copied %v bytes", bytesCopied)
	}

	responseMsg, err := json.Marshal(msg)

	if err != nil {
		log.Println("[ERROR] Couldn't marshal json", err)

		response.WriteHeader(http.StatusInternalServerError)
		msg.Message = "Something went wrong"
	}

	response.Write(responseMsg)
}
