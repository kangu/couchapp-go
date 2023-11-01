package main

import (
	"couchappp-go/tests"
	"encoding/base64"
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
