package circleci

const projectEnvVarPath = "/envvar"

// ProjectEnvVarService is an interface for ProjectEnvVar in Project API.
type ProjectEnvVarService interface {
	Create(projectSlug, name, value string) (*ProjectEnvVar, error)
	Get(projectSlug, name string) (*ProjectEnvVar, error)
	List(projectSlug string) (*ProjectEnvVarList, error)
	Delete(projectSlug, name string) error
}

// ProjectEnvVarOp handles communication with the project related methods in the CircleCI API v2.
type ProjectEnvVarOp struct {
	client *Client
}

var _ ProjectEnvVarService = (*ProjectEnvVarOp)(nil)

// ProjectEnvVar represents an environment variable in a Project.
type ProjectEnvVar struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// ProjectEnvVar represents a list of ProjectEnvVar.
type ProjectEnvVarList struct {
	NextPageToken string           `json:"next_page_token,omitempty"`
	Items         []*ProjectEnvVar `json:"items,omitempty"`
}

// Create adds a new environment variable or update existing variable on the specified project.
// Returns the added env var (the value will be masked).
func (ps *ProjectEnvVarOp) Create(projectSlug, name, value string) (*ProjectEnvVar, error) {
	ev := &ProjectEnvVar{}
	err := ps.client.Post(envVarPathPrefix(projectSlug), &ProjectEnvVar{Name: name, Value: value}, ev)
	if err != nil {
		return nil, err
	}
	return ev, nil
}

// Get gets environment variable.
// Returns the env vars (the value will be masked).
func (ps *ProjectEnvVarOp) Get(projectSlug, name string) (*ProjectEnvVar, error) {
	ev := &ProjectEnvVar{}
	err := ps.client.Get(envVarValuePathPrefix(projectSlug, name), ev, nil)
	if err != nil {
		return nil, err
	}
	return ev, nil
}

// List list environment variable to the specified project.
// Returns the env vars (the value will be masked).
func (ps *ProjectEnvVarOp) List(projectSlug string) (*ProjectEnvVarList, error) {
	evp := &ProjectEnvVarList{}
	err := ps.client.Get(envVarPathPrefix(projectSlug), evp, nil)
	if err != nil {
		return nil, err
	}
	return evp, nil
}

// Delete deletes the specified environment variable from the project.
func (ps *ProjectEnvVarOp) Delete(projectSlug, name string) error {
	return ps.client.Delete(envVarValuePathPrefix(projectSlug, name))
}

func envVarPathPrefix(projectSlug string) string {
	return projectPathPrefix(projectSlug) + projectEnvVarPath
}
func envVarValuePathPrefix(projectSlug, envVar string) string {
	return envVarPathPrefix(projectSlug) + "/" + envVar
}
