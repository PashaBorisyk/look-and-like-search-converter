package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

const contentType string = "application/json"
const apiKey string = "FEBA1843F8CAA5916F0465CDD3F4A256"
const url = "https://lookandlike.search.windows.net/indexes/look-and-like-test/docs/index?api-version=2019-05-06"

func UploadModelToIndex(model interface{}) error {
	bodyContent, err := json.Marshal(model)
	log.Println(string(bodyContent))
	request, err := http.NewRequest("POST", url, bytes.NewReader(bodyContent))
	if err != nil {
		log.Println("Error while creating a document: ", err)
		return err
	}
	request.Header.Add("Content-Type", contentType)
	request.Header.Add("api-key", apiKey)
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println("Error performing POST request: ", err)
		return err
	}
	if !(resp.StatusCode == 200 || resp.StatusCode == 201 || resp.StatusCode == 202) {
		log.Println("Not successful post request. Exiting with server status code : ", resp.StatusCode)
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(resp.Body)
		newStr := buf.String()
		log.Println("Resp body:", newStr)
		log.Println("Request body was: ",string(bodyContent))
		return errors.New("Request returned not successful status :" + resp.Status)
	}
	return nil
}