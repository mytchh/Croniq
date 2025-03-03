package handlers

import (
	"context"
	"croniq/backend/k8s"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateCronJobRequest struct {
	Name      string   `json:"name"`
	Namespace string   `json:"namespace"`
	Schedule  string   `json:"schedule"`
	Image     string   `json:"image"`
	Command   []string `json:"command,omitempty"`
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
	case http.MethodPost:
		h.handleCreateCronJob(w, r)
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

func (h *CronJobHandler) handleCreateCronJob(w http.ResponseWriter, r *http.Request) {
	var req CreateCronJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	client, err := h.k8sClient.GetClient()
	if err != nil {
		sendJSONError(w, fmt.Sprintf("Failed to create K8s client: %v", err), http.StatusInternalServerError)
		return
	}

	cronJob := &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
		Spec: batchv1.CronJobSpec{
			Schedule: req.Schedule,
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							RestartPolicy: corev1.RestartPolicyOnFailure,
							Containers: []corev1.Container{
								{
									Name:    req.Name,
									Image:   req.Image,
									Command: req.Command,
								},
							},
						},
					},
				},
			},
		},
	}

	result, err := client.BatchV1().CronJobs(req.Namespace).Create(r.Context(), cronJob, metav1.CreateOptions{})
	if err != nil {
		sendJSONError(w, fmt.Sprintf("Failed to create cron job: %v", err), http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, result)
}

func sendJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		sendJSONError(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
