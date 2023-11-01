package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
)

// doCouchRequest runs an HTTP request to the couch
// only used internally by other functions
func doCouchRequest(url string, method string, auth string, body []byte) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return &http.Response{}, err
	}

	// all requests have the json type
	req.Header.Set("Content-Type", "application/json")
	// Add Basic Authentication header
	if auth != "" {
		req.Header.Add("Authorization", auth)
	}

	resp, err := client.Do(req)
	if err != nil {
		return &http.Response{}, err
	}

	return resp, nil
}

func deleteDatabase(host string, db string, auth string) (bool, error) {
	resp, err := doCouchRequest(fmt.Sprintf("%s/%s", host, db), "DELETE", auth, nil)
	if err != nil {
		return false, err
	}
	return resp.StatusCode == 200, nil
}

func getDocRevision(host string, db string, docid string, auth string) (string, error) {
	// first try to create database
	// should fail silently if already present
	respCreation, err := doCouchRequest(fmt.Sprintf("%s/%s", host, db), "PUT", auth, nil)
	if err != nil {
		return "", err
	}
	defer respCreation.Body.Close()
	if respCreation.StatusCode == http.StatusUnauthorized {
		return "", errors.New(fmt.Sprintf("Error setting up db: %s", respCreation.Status))
	}

	// look for existing doc
	url := fmt.Sprintf("%s/%s/%s", host, db, docid)
	resp, err := doCouchRequest(url, "HEAD", auth, nil)
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
	resp, err := doCouchRequest(url, "POST", auth, doc)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return 0, errors.New(fmt.Sprintf("Error uploading document: %s", resp.Status))
	}

	return resp.StatusCode, nil
}

func getDoc(host string, db string, id string, auth string) (map[string]interface{}, error) {
	var result map[string]interface{}
	url := fmt.Sprintf("%s/%s/%s", host, db, id)
	resp, err := doCouchRequest(url, "GET", auth, nil)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}
	if err = json.Unmarshal(body, &result); err != nil {
		return result, err
	}
	return result, nil
}

func checkIdenticalDocs(host string, db string, doc map[string]interface{}, auth string) (bool, error) {
	idStr, ok := doc["_id"].(string)
	if !ok {
		return false, errors.New("doc id should be a string")
	}
	obj2, err := getDoc(host, db, idStr, auth)
	if err != nil {
		return false, err
	}
	// take out the revision as that doesn't need comparing
	delete(obj2, "_rev")

	if reflect.DeepEqual(doc, obj2) {
		return true, nil
	}
	return false, nil
}
