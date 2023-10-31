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

func getDocRevision(host string, db string, docid string, auth string) (string, error) {
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

func checkIdenticalDocs(host string, db string, doc map[string]interface{}, auth string) (bool, error) {
	url := fmt.Sprintf("%s/%s/%s", host, db, doc["_id"])
	resp, err := doCouchRequest(url, "GET", auth, nil)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// Handle the error
		fmt.Println("Error reading response body:", err)
		return false, err
	}

	var obj2 map[string]interface{}
	if err := json.Unmarshal(body, &obj2); err != nil {
		fmt.Println("Error unmarshaling existing json doc:", err)
		return false, err
	}
	// take out the revision as that doesn't need comparing
	delete(obj2, "_rev")

	if reflect.DeepEqual(doc, obj2) {
		return true, nil
	}
	return false, nil
}
