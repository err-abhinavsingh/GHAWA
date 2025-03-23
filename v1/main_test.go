package main

import (
	"GHAWA/entities"
	"GHAWA/server"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestRiskServer(t *testing.T) {
	go func() {
		server.StartServer()
	}()

	t.Run("", testGetRisks)
	t.Run("", testGetRisk)
	t.Run("", testCreateRisk)
}

/*
This test tests:
1. Getting risks when none are present
2. Getting risks when some are present
*/
func testGetRisks(t *testing.T) {

	resp, err := http.Get("http://localhost:8080/v1/risks")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected OK, got %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var risks []*entities.Risk
	err = json.Unmarshal(body, &risks)
	if err != nil {
		t.Fatal(err)
	}

	if len(risks) != 0 {
		t.Errorf("Expected 0 risks, got %v", len(risks))
	}

	risk := &entities.Risk{
		Description: "Risk Number 1",
		Title:       "Risk1",
		State:       "open",
	}
	body, err = json.Marshal(risk)
	if err != nil {
		t.Fatal(err)
	}
	_, err = http.Post("http://localhost:8080/v1/risks", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	risk = &entities.Risk{
		Description: "Risk Number 2",
		Title:       "Risk2",
		State:       "investigating",
	}
	body, err = json.Marshal(risk)
	if err != nil {
		t.Fatal(err)
	}
	_, err = http.Post("http://localhost:8080/v1/risks", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	resp, err = http.Get("http://localhost:8080/v1/risks")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected OK, got %v", resp.StatusCode)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	risks = nil
	err = json.Unmarshal(body, &risks)
	if err != nil {
		t.Fatal(err)
	}

	if len(risks) != 2 {
		t.Errorf("Expected 2 risks, got %v", len(risks))
	}
}

/*
This test tests:
1. Getting a risk for a uuid which is not present
2. Getting a risk for a uuid which is present
*/
func testGetRisk(t *testing.T) {

	resp, err := http.Get("http://localhost:8080/v1/risks/1234")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected OK, got %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var risk *entities.Risk
	err = json.Unmarshal(body, &risk)
	if err != nil {
		t.Fatal(err)
	}

	if risk != nil {
		t.Errorf("Expected no risk, got %v", risk)
	}

	risk = &entities.Risk{
		Description: "Risk Number 1",
		Title:       "Risk1",
		State:       "open",
	}
	body, err = json.Marshal(risk)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.Post("http://localhost:8080/v1/risks", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	risk = nil
	err = json.Unmarshal(body, &risk)
	if err != nil {
		t.Fatal(err)
	}

	expectedUuid := risk.Uuid

	resp, err = http.Get(fmt.Sprintf("http://localhost:8080/v1/risks/%s", expectedUuid))
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected OK, got %v", resp.StatusCode)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	risk = nil
	err = json.Unmarshal(body, &risk)
	if err != nil {
		t.Fatal(err)
	}

	if risk == nil {
		t.Errorf("Expected 1 risk, got nil")
	}
	if risk.Uuid != expectedUuid {
		t.Errorf("Expected %s, got %s", expectedUuid, risk.Uuid)
	}
}

/*
This test tests creation of a Risk with Invalid State.
*/
func testCreateRisk(t *testing.T) {

	risk := &entities.Risk{
		Description: "Risk Number 1",
		Title:       "Risk1",
		State:       "invalid state",
	}
	body, err := json.Marshal(risk)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post("http://localhost:8080/v1/risks", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected BadRequest, got %v", resp.StatusCode)
	}
}
