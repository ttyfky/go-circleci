// Project API of CircleCI
//https://circleci.com/docs/api/v2/#tag/Project

package circleci

const projectBasePath = "/project"

// ProjectService is an interface for Project API.
type ProjectService interface {
	Get(projectSlug string) (*Project, error)
}

// ProjectServiceOp handles communication with the project related methods in the CircleCI API v2.
type ProjectServiceOp struct {
	client *Client
}

// Project represents information about a project in CircleCI.
type Project struct {
	Slug             string `json:"slug,omitempty"`
	Name             string `json:"name,omitempty"`
	OrganizationName string `json:"organization_name,omitempty"`
	ExternalURL      string `json:"external_url,omitempty"`
	VcsInfo          struct {
		VcsURL        string `json:"vcs_url,omitempty"`
		Provider      string `json:"provider,omitempty"`
		DefaultBranch string `json:"default_branch,omitempty"`
	} `json:"vcs_info,omitempty"`
}

// Get gets project information.
func (ps *ProjectServiceOp) Get(projectSlug string) (*Project, error) {
	p := &Project{}
	err := ps.client.Get(projectPathPrefix(projectSlug), p, nil)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func projectPathPrefix(projectSlug string) string {
	return projectBasePath + "/" + projectSlug
}
