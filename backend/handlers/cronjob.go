package handlers

import (
	"context"
	"croniq/backend/k8s"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type CronJobHandler struct {
	k8sClient *k8s.Client
}

func NewCronJobHandler(client *k8s.Client) *CronJobHandler {
	return &CronJobHandler{
		k8sClient: client,
	}
}

func (h *CronJobHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetCronJobs(w, r)
	default:
		sendJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *CronJobHandler) handleGetCronJobs(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	client, err := h.k8sClient.GetClient()
	if err != nil {
		sendJSONError(w, fmt.Sprintf("Failed to create K8s client: %v", err), http.StatusInternalServerError)
		return
	}

	cronJobs, err := client.BatchV1().CronJobs("").List(ctx, metav1.ListOptions{})
	if err != nil {
		sendJSONError(w, fmt.Sprintf("Failed to fetch cron jobs: %v", err), http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, cronJobs)
}

func sendJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		sendJSONError(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
