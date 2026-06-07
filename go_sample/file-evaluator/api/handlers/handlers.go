package handlers

import (
	"encoding/json"
	producer "file-evaluator/api/kafka"
	database "file-evaluator/db"
	"fmt"
	"net/http"
)

type Payload struct {
	FilePath string `json:"path"`
	EvalType string `json:"type"`
}

type HttpHandler struct {
	DB            *database.Store
	KafkaProducer *producer.KafkaProducer
}

type Message struct {
	ID       int64  `json:"id"`
	FilePath string `json:"path"`
	EvalType string `json:"type"`
}

func (h *HttpHandler) FileEvalHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodPost:
		h.handlePost(w, r)

	case http.MethodGet:
		h.handleGet(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *HttpHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	payload, payloadErr := validatePostPayload(r)
	if payloadErr != nil {
		http.Error(w, payloadErr.Error(), http.StatusBadRequest)
	}
	requestID, DBErr := h.DB.CreateNewEvent(payload.EvalType)
	if DBErr != nil {
		http.Error(w, DBErr.Error(), http.StatusInternalServerError)
	}

	kafkaErr := produceKafkaMessage(h.KafkaProducer, requestID, payload)
	if kafkaErr != nil {
		http.Error(w, kafkaErr.Error(), http.StatusInternalServerError)
	}

	response := struct {
		ID int64 `json:"id"`
	}{
		ID: requestID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *HttpHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	id, paramErr := validateGetParams(r)
	if paramErr != nil {
		http.Error(w, paramErr.Error(), http.StatusBadRequest)
		return
	}

	result, DBErr := h.DB.GetEvent(id)
	if DBErr != nil {
		http.Error(w, DBErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func validateGetParams(r *http.Request) (int64, error) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		return 0, fmt.Errorf("Missing id parameter: id")
	}

	var id int64
	_, err := fmt.Sscanf(idParam, "%d", &id)
	if err != nil {
		return 0, fmt.Errorf("Id must be an int")
	}

	return id, nil
}

func getEventFromDatabase(id int64, DB *database.Store) ([]database.EvaluationResult, error) {
	result, err := DB.GetEvent(id)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch event: %v", err)
	}

	if result == nil {
		return nil, fmt.Errorf("ID not found: %d", id)
	}

	return result, nil

}

func validatePostPayload(r *http.Request) (Payload, error) {
	var payload Payload

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		return Payload{}, fmt.Errorf("Invalid JSON body: %v", err)
	}

	if payload.FilePath == "" {
		return Payload{}, fmt.Errorf("Missing required field: path")
	}
	return payload, nil
}

func produceKafkaMessage(producer *producer.KafkaProducer, requestID int64, payload Payload) error {
	kafkaMessageWithID := Message{
		ID:       requestID,
		FilePath: payload.FilePath,
		EvalType: payload.EvalType,
	}

	kafkaMessageWithIDBytes, err := json.Marshal(kafkaMessageWithID)
	if err != nil {
		return fmt.Errorf("Failed to formulate Kafka message: %v", err)
	}

	var kafkaErr error
	kafkaErr = producer.ProduceEvent(kafkaMessageWithIDBytes)
	if kafkaErr != nil {
		return fmt.Errorf("Failed to produce Kafka message: %v", err)
	}
	return nil
}
