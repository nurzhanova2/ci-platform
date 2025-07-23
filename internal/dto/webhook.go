package dto

type WebhookPayload struct {
	Repository struct {
		CloneURL string `json:"clone_url"`
		Name     string `json:"name"`
	} `json:"repository"`
	Ref string `json:"ref"` 
}
