package main

import (
	"croniq/backend/handlers"
	"croniq/backend/k8s"
	"log"
	"net/http"
	"os"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	kubeconfigPath := os.Getenv("KUBECONFIG")
	k8sClient := k8s.NewClient(kubeconfigPath)
	cronJobHandler := handlers.NewCronJobHandler(k8sClient)

	mux := http.NewServeMux()
	mux.Handle("/api/cronjobs", corsMiddleware(cronJobHandler))
	mux.HandleFunc("/api/cluster-info", corsMiddleware(http.HandlerFunc(cronJobHandler.HandleClusterInfo)).ServeHTTP)
	mux.HandleFunc("/api/jobs", corsMiddleware(http.HandlerFunc(cronJobHandler.HandleJobs)).ServeHTTP)
	mux.HandleFunc("/api/stats", corsMiddleware(http.HandlerFunc(cronJobHandler.HandleStats)).ServeHTTP)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Printf("Server starting on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
