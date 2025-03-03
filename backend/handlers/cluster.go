package handlers

import (
	"net/http"

	"k8s.io/client-go/tools/clientcmd"
)

type ClusterInfo struct {
	Name          string `json:"name"`
	ServerAddress string `json:"serverAddress"`
	Version       string `json:"version"`
}

func (h *CronJobHandler) HandleClusterInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	client, err := h.k8sClient.GetClient()
	if err != nil {
		sendJSONError(w, "Failed to get cluster info", http.StatusInternalServerError)
		return
	}

	version, err := client.Discovery().ServerVersion()
	if err != nil {
		sendJSONError(w, "Failed to get server version", http.StatusInternalServerError)
		return
	}

	config, err := h.k8sClient.GetConfig()
	if err != nil {
		sendJSONError(w, "Failed to get cluster config", http.StatusInternalServerError)
		return
	}

	// Get the current context name
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	rawConfig, err := kubeConfig.RawConfig()
	if err != nil {
		sendJSONError(w, "Failed to get raw config", http.StatusInternalServerError)
		return
	}

	info := ClusterInfo{
		Name:          rawConfig.CurrentContext,
		ServerAddress: config.Host,
		Version:       version.String(),
	}

	sendJSONResponse(w, info)
}
