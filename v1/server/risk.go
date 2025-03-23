package server

import (
	"GHAWA/entities"
	"GHAWA/services"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	riskService *services.RiskService
)

func StartServer() {
	riskService = services.NewRiskService()

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/risks/{id}", riskHandler)
	mux.HandleFunc("/v1/risks", risksHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

		<-signalChannel

		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	log.Fatal(server.ListenAndServe())
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
		http.Error(w, "", http.StatusMethodNotAllowed)
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
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}
