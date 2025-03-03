package handlers

import (
	"context"
	"net/http"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type JobStats struct {
	TotalCronJobs  int `json:"totalCronJobs"`
	ActiveCronJobs int `json:"activeCronJobs"`
	TotalJobs      int `json:"totalJobs"`
	RunningJobs    int `json:"runningJobs"`
	FailedJobs     int `json:"failedJobs"`
	SucceededJobs  int `json:"succeededJobs"`
}

func (h *CronJobHandler) HandleJobs(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	client, err := h.k8sClient.GetClient()
	if err != nil {
		sendJSONError(w, "Failed to create K8s client", http.StatusInternalServerError)
		return
	}

	jobs, err := client.BatchV1().Jobs("").List(ctx, metav1.ListOptions{})
	if err != nil {
		sendJSONError(w, "Failed to list jobs", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, jobs)
}

func (h *CronJobHandler) HandleStats(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	client, err := h.k8sClient.GetClient()
	if err != nil {
		sendJSONError(w, "Failed to create K8s client", http.StatusInternalServerError)
		return
	}

	cronJobs, err := client.BatchV1().CronJobs("").List(ctx, metav1.ListOptions{})
	if err != nil {
		sendJSONError(w, "Failed to list cronjobs", http.StatusInternalServerError)
		return
	}

	jobs, err := client.BatchV1().Jobs("").List(ctx, metav1.ListOptions{})
	if err != nil {
		sendJSONError(w, "Failed to list jobs", http.StatusInternalServerError)
		return
	}

	stats := JobStats{
		TotalCronJobs:  len(cronJobs.Items),
		ActiveCronJobs: 0,
		TotalJobs:      len(jobs.Items),
		RunningJobs:    0,
		FailedJobs:     0,
		SucceededJobs:  0,
	}

	for _, cj := range cronJobs.Items {
		if cj.Spec.Suspend == nil || !*cj.Spec.Suspend {
			stats.ActiveCronJobs++
		}
	}

	for _, job := range jobs.Items {
		if job.Status.Active > 0 {
			stats.RunningJobs++
		}
		if job.Status.Failed > 0 {
			stats.FailedJobs++
		}
		if job.Status.Succeeded > 0 {
			stats.SucceededJobs++
		}
	}

	sendJSONResponse(w, stats)
}
