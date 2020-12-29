package circleci

import "time"

const (
	contextBasePath = "/context"

	defaultContextType = "organization"
)

// ContextService is an interface for Context in Project API.
type ContextService interface {
	List(slug string) (*ContextList, error)
	Create(projectSlug, name string) (*Context, error)
	Delete(id string) error
	Get(id string) (*Context, error)
	ListEnvVar(id string) (*ContextEnvVarList, error)
	UpsertEnvVar(id, envVarName, envVarValue string) (*ContextEnvVar, error)
	RemoveEnvVar(id, envVarName string) error
}

// ContextOp handles communication with the project related methods in the CircleCI API v2.
type ContextOp struct {
	client *Client
}

var _ ContextService = (*ContextOp)(nil)

// Context represents information about a context in CircleCI.
type Context struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// ContextList represents a list of Context variables.
type ContextList struct {
	NextPageToken string     `json:"next_page_token,omitempty"`
	Items         []*Context `json:"items,omitempty"`
}

// Owner represents Owner used in ContextCreate.
type Owner struct {
	Slug string `json:"slug,omitempty"`
	Type string `json:"type,omitempty"`
}

// ContextCreate represents payload to create Context.
type ContextCreate struct {
	Name   string `json:"name,omitempty"`
	*Owner `json:"owner,omitempty"`
}

// ContextEnvVar represents information about an environment variable of Context.
type ContextEnvVar struct {
	Variable  string    `json:"variable,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	ContextID string    `json:"context_id,omitempty"`
}

// ContextEnvVarList represents a list of ContextEnvVar.
type ContextEnvVarList struct {
	NextPageToken string           `json:"next_page_token,omitempty"`
	Items         []*ContextEnvVar `json:"items,omitempty"`
}

// List list contexts for an owner.
// owner-slug is expected but not owner-id.
func (ps *ContextOp) List(slug string) (*ContextList, error) {
	cl := &ContextList{}
	path := contextBasePath + "?owner-slug=" + slug
	err := ps.client.Get(path, cl, nil)
	if err != nil {
		return nil, err
	}
	return cl, nil
}

// Create adds a new environment variable or update existing variable on the specified project.
// Returns the added env var (the value will be masked).
func (ps *ContextOp) Create(projectSlug, name string) (*Context, error) {
	c := &Context{}
	err := ps.client.Post(contextBasePath, &ContextCreate{Name: name, Owner: &Owner{Slug: projectSlug, Type: defaultContextType}}, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Delete deletes the specified environment variable from the project.
func (ps *ContextOp) Delete(id string) error {
	return ps.client.Delete(contextBasePath + "/" + id)
}

// Get gets environment variable.
// Returns the env vars (the value will be masked).
func (ps *ContextOp) Get(id string) (*Context, error) {
	c := &Context{}
	err := ps.client.Get(contextBasePath+"/"+id, c, nil)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// ListEnvVar list contexts for an owner.
// Returns the env vars (the value will be masked).
func (ps *ContextOp) ListEnvVar(id string) (*ContextEnvVarList, error) {
	cel := &ContextEnvVarList{}
	err := ps.client.Get(contextEnvVarPath(id), cel, nil)
	if err != nil {
		return nil, err
	}
	return cel, nil
}

// UpsertEnvVar list contexts for an owner.
// Returns the env vars (the value will be masked).
func (ps *ContextOp) UpsertEnvVar(id, envVarName, envVarValue string) (*ContextEnvVar, error) {
	ce := &ContextEnvVar{}
	path := contextEnvVarPath(id) + "/" + envVarName
	err := ps.client.Put(path,
		struct {
			Value string `json:"value"`
		}{
			Value: envVarValue,
		},
		ce)
	if err != nil {
		return nil, err
	}
	return ce, nil
}

// RemoveEnvVar list contexts for an owner.
// Returns the env vars (the value will be masked).
func (ps *ContextOp) RemoveEnvVar(id, envVarName string) error {
	return ps.client.Delete(contextEnvVarPath(id) + "/" + envVarName)
}

func contextEnvVarPath(id string) string {
	return contextBasePath + "/" + id + "/environment-variable"
}
