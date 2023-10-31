package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func getDocRevision(host string, db string, docid string, auth string) (string, error) {
	url := fmt.Sprintf("%s/%s/%s", host, db, docid)

	client := &http.Client{}

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return "", err
	}

	// Add Basic Authentication header
	if auth != "" {
		req.Header.Add("Authorization", auth)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusNotModified:
		{
			unquotedHeaderValue, err := strconv.Unquote(resp.Header["Etag"][0])
			if err != nil {
				return "", err
			}
			return unquotedHeaderValue, nil
		}
	case http.StatusUnauthorized:
		return "", errors.New(fmt.Sprintf("Error getting revision: %s", resp.Status))
	case http.StatusNotFound:
		return "", nil
	default:
		// default to error if expected status is not returned
		return "", errors.New(fmt.Sprintf("Error getting revision: %s", resp.Status))
	}

}

// postDoc pushed the document to the CouchDB and returns the http status code
func postDoc(host string, db string, doc []byte, auth string) (int, error) {
	url := fmt.Sprintf("%s/%s", host, db)

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(doc))
	if err != nil {
		return 0, err
	}
	// set proper headers
	req.Header.Set("Content-Type", "application/json")
	if auth != "" {
		req.Header.Add("Authorization", auth)
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return 0, errors.New(fmt.Sprintf("Error uploading document: %s", resp.Status))
	}

	return resp.StatusCode, nil
}
