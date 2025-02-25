package main

import (
	"bytes"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

func Post(reqUrl string, auth string, body interface{}) (respBodyObj ResponseBody, err error) {
	postBody, _ := json.Marshal(body)
	requestBody := bytes.NewBuffer(postBody)
	log.WithFields(log.Fields{
		"body": string(postBody),
	}).Debug("The request body")
	req, err := http.NewRequest("POST", reqUrl, requestBody)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", auth)
	return handleResp(req)
}

func Get(reqUrl string, auth string) (respBodyObj ResponseBody, err error) {
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", auth)
	return handleResp(req)
}

func handleResp(req *http.Request) (respBodyObj ResponseBody, err error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	log.WithFields(log.Fields{
		"body": string(respBody),
	}).Debug("The response body")
	if resp.StatusCode != 200 {
		return respBodyObj, errors.New("received non 200 response code. The response code was " + strconv.Itoa(resp.StatusCode))
	}

	err = json.Unmarshal(respBody, &respBodyObj)
	if err != nil {
		log.Fatalln("There was error while parsing the response from server. Exiting...")
	}

	return respBodyObj, nil
}
