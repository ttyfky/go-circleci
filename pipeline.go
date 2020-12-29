package circleci

import "time"

// Pipeline represents pipeline in CircleCI.
type Pipeline struct {
	ID     string `json:"id,omitempty"`
	Errors []struct {
		Type    string `json:"type,omitempty"`
		Message string `json:"message,omitempty"`
	} `json:"errors,omitempty"`
	ProjectSlug string    `json:"project_slug,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	Number      int       `json:"number,omitempty"`
	State       string    `json:"state,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	Trigger     struct {
		Type       string    `json:"type,omitempty"`
		ReceivedAt time.Time `json:"received_at,omitempty"`
		Actor      struct {
			Login     string `json:"login,omitempty"`
			AvatarURL string `json:"avatar_url,omitempty"`
		} `json:"actor,omitempty"`
	} `json:"trigger,omitempty"`
	Vcs struct {
		ProviderName        string `json:"provider_name,omitempty"`
		TargetRepositoryURL string `json:"target_repository_url,omitempty"`
		Branch              string `json:"branch,omitempty"`
		ReviewID            string `json:"review_id,omitempty"`
		ReviewURL           string `json:"review_url,omitempty"`
		Revision            string `json:"revision,omitempty"`
		Tag                 string `json:"tag,omitempty"`
		Commit              struct {
			Subject string `json:"subject,omitempty"`
			Body    string `json:"body"`
		} `json:"commit"`
		OriginRepositoryURL string `json:"origin_repository_url"`
	} `json:"vcs"`
}
