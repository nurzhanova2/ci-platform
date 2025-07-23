package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/nurzhanova2/ci-platform/internal/service"
)

type WebhookHandler struct {
	PipelineService *service.PipelineService
}

type webhookPayload struct {
	RepositoryURL string `json:"repository_url"`
	Branch        string `json:"branch"`
	CommitSHA     string `json:"commit_sha"`
}

func NewWebhookHandler(pipelineService *service.PipelineService) *WebhookHandler {
	return &WebhookHandler{PipelineService: pipelineService}
}

func (h *WebhookHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var payload webhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// валидация обязательного поля
	if payload.RepositoryURL == "" {
		http.Error(w, "missing repository_url", http.StatusBadRequest)
		return
	}

	go func() {
		err := h.PipelineService.RunPipeline(payload.RepositoryURL, payload.Branch)
		if err != nil {
			log.Printf("pipeline error: %v", err)
		}
	}()

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("pipeline started"))
}
