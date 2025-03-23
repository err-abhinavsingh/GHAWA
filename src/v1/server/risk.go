package main

import (
	"GHAWA/entities"
	"GHAWA/services"
	"encoding/json"
	"log"
	"net/http"
)

var (
	riskService *services.RiskService
)

func main() {
	riskService = services.NewRiskService()

	http.HandleFunc("/v1/risks/{id}", riskHandler)
	http.HandleFunc("/v1/risks", risksHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func riskHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Print("ERROR: ", "recovered in riskHandler ", r)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
	}()

	switch r.Method {
	case http.MethodGet:
		riskId := r.PathValue("id")
		var risk *entities.Risk
		risk = riskService.GetRisk(riskId)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(risk)
		if err != nil {
			log.Print("ERROR: ", "failed to encode risk ", err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func risksHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Print("ERROR: ", "recovered in risksHandler ", r)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
	}()

	switch r.Method {
	case http.MethodGet:
		var risks []*entities.Risk
		risks = riskService.GetRisks()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(risks)
		if err != nil {
			log.Print("ERROR: ", "failed to encode risks ", err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
	case http.MethodPost:
		var risk = entities.Risk{}
		err := json.NewDecoder(r.Body).Decode(&risk)
		if err != nil {
			log.Print("ERROR: ", "failed to decode risk ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var createdRisk *entities.Risk
		createdRisk, err = riskService.CreateRisk(risk)
		if err != nil {
			log.Print("ERROR: ", "failed to create risk ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(createdRisk)
		if err != nil {
			log.Print("ERROR: ", "failed to encode risk ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
