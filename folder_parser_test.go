package main

import (
	"couchappp-go/tests"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestFolderToJSONS_SimpleDesign(t *testing.T) {
	result, err := FolderToJSON("tests/fixtures/simple_design")
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	expected := "_design/simple"
	if result["_id"] != expected {
		t.Errorf("Expected %s, but got %s", expected, result["_id"])
	}
}

func TestProcessFolderToCouch_SimpleDesign(t *testing.T) {
	const TargetDB = "test1"

	authTestConfig, err := tests.LoadTestConfig()
	if err != nil {
		t.Errorf("Config file not properly setup %v", err)
	}
	params := CliParams{
		source: "tests/fixtures/simple_design",
		db:     TargetDB,
		host:   "http://localhost:5984",
	}
	if authTestConfig.Username != "" {
		params.base64auth = "Basic " + base64.StdEncoding.EncodeToString([]byte(authTestConfig.Username+":"+authTestConfig.Password))
	}

	// cleanup before pushing
	_, _ = deleteDatabase(params.host, params.db, params.base64auth)

	ProcessFolderToCouch(params)

	// fetch document and make sure id matches
	expected := "_design/simple"
	doc, _ := getDoc(params.host, params.db, expected, params.base64auth)
	if doc["_id"] != expected {
		t.Errorf("Expected %s, but got %s", expected, doc["_id"])
	}
}

func TestProcessFolderToCouch_WithReduce(t *testing.T) {
	const TargetDB = "test1"

	authTestConfig, err := tests.LoadTestConfig()
	if err != nil {
		t.Errorf("Config file not properly setup %v", err)
	}
	params := CliParams{
		source: "tests/fixtures/with_reduce",
		db:     TargetDB,
		host:   "http://localhost:5984",
	}
	if authTestConfig.Username != "" {
		params.base64auth = "Basic " + base64.StdEncoding.EncodeToString([]byte(authTestConfig.Username+":"+authTestConfig.Password))
	}

	// cleanup before pushing
	_, _ = deleteDatabase(params.host, params.db, params.base64auth)

	ProcessFolderToCouch(params)

	// post some documents to query
	doc1 := make(map[string]interface{})
	doc1["_id"] = "test_doc"
	// Convert the doc to a JSON string
	jsonString, err := json.MarshalIndent(doc1, "", "    ")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	_, err = postDoc(params.host, params.db, jsonString, params.base64auth)
	if err != nil {
		fmt.Println("Error uploading document:", err)
		return
	}

	// another document
	doc2 := map[string]interface{}{
		"_id": "test_doc_2",
	}
	jsonString, err = json.MarshalIndent(doc2, "", "    ")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	_, err = postDoc(params.host, params.db, jsonString, params.base64auth)
	if err != nil {
		fmt.Println("Error uploading document:", err)
		return
	}

	// fetch view and check for reduce result
	doc, _ := getView(params.host, params.db, "simple", "test", params.base64auth)
	expected := float64(2)
	if doc.Rows[0].Value != expected {
		t.Errorf("Expected %v, but got %v", expected, doc.Rows[0].Value)
	}
}
