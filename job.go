package circleci

import "time"

const jobBasePath = "/job"

type JobService interface {
	Get(id, projectSlug string) (*Job, error)
	Cancel(id, projectSlug string) (*Message, error)
	GetArtifacts(id, projectSlug string) (*ArtifactList, error)
	GetTestMetadata(id, projectSlug string) (*TestMetadataList, error)
}

// JobOp handles communication with the project related methods in the CircleCI API v2.
type JobOp struct {
	client *Client
}

var _ JobService = (*JobOp)(nil)

// Job represents information about job in CircleCI.
type Job struct {
	WebURL       string `json:"web_url,omitempty"`
	Project      `json:"project,omitempty"`
	ParallelRuns []struct {
		Index  int    `json:"index,omitempty"`
		Status string `json:"status,omitempty"`
	} `json:"parallel_runs,omitempty"`
	StartedAt      time.Time `json:"started_at,omitempty"`
	LatestWorkflow struct {
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"latest_workflow,omitempty"`
	Name     string `json:"name,omitempty"`
	Executor struct {
		Type          string `json:"type,omitempty"`
		ResourceClass string `json:"resource_class,omitempty"`
	} `json:"executor,omitempty"`
	Parallelism  int         `json:"parallelism,omitempty"`
	Status       interface{} `json:"status,omitempty"`
	Number       int         `json:"number,omitempty"`
	Pipeline     Pipeline    `json:"pipeline,omitempty"`
	Duration     int         `json:"duration,omitempty"`
	CreatedAt    time.Time   `json:"created_at,omitempty"`
	Messages     []Message   `json:"messages,omitempty"`
	Contexts     []Context   `json:"contexts,omitempty"`
	Organization struct {
		Name string `json:"name,omitempty"`
	} `json:"organization,omitempty"`
	QueuedAt  time.Time `json:"queued_at,omitempty"`
	StoppedAt time.Time `json:"stopped_at,omitempty"`
}

// Jobs represents information about a jobs having dependencies.
type Jobs struct {
	Items []struct {
		CanceledBy        string      `json:"canceled_by,omitempty"`
		Dependencies      []string    `json:"dependencies,omitempty"`
		JobNumber         int         `json:"job_number,omitempty"`
		ID                string      `json:"id,omitempty"`
		StartedAt         time.Time   `json:"started_at,omitempty"`
		Name              string      `json:"name,omitempty"`
		ApprovedBy        string      `json:"approved_by,omitempty"`
		ProjectSlug       string      `json:"project_slug,omitempty"`
		Status            interface{} `json:"status,omitempty"`
		Type              string      `json:"type,omitempty"`
		StoppedAt         time.Time   `json:"stopped_at,omitempty"`
		ApprovalRequestID string      `json:"approval_request_id,omitempty"`
	} `json:"items,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
}

// Metadata represents information about a test metadata of a Job.
type Metadata struct {
	Message   string `json:"message,omitempty"`
	Source    string `json:"source,omitempty"`
	RunTime   string `json:"run_time,omitempty"`
	File      string `json:"file,omitempty"`
	Result    string `json:"result,omitempty"`
	Name      string `json:"name,omitempty"`
	Classname string `json:"classname,omitempty"`
}

// TestMetadataList contains list of Metadata of test.
type TestMetadataList struct {
	Items         []Metadata `json:"items,omitempty"`
	NextPageToken string     `json:"next_page_token,omitempty"`
}

// Artifact represents artifact of a Job.
type Artifact struct {
	Path      string `json:"path,omitempty"`
	NodeIndex int    `json:"node_index,omitempty"`
	URL       string `json:"url,omitempty"`
}

// ArtifactList represents list of Artifact in a Job.
type ArtifactList struct {
	Items         []Artifact `json:"items,omitempty"`
	NextPageToken string     `json:"next_page_token,omitempty"`
}

// Get gets job detail.
func (ps *JobOp) Get(id, projectSlug string) (*Job, error) {
	j := &Job{}
	path := jobIDPath(id, projectSlug)
	err := ps.client.Get(path, j, nil)
	if err != nil {
		return nil, err
	}
	return j, nil
}

// Cancel cancels a given job.
func (ps *JobOp) Cancel(id, projectSlug string) (*Message, error) {
	m := &Message{}
	path := jobIDPath(id, projectSlug) + "/cancel"
	err := ps.client.Post(path, nil, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// GetArtifacts get artifacts in a job.
func (ps *JobOp) GetArtifacts(id, projectSlug string) (*ArtifactList, error) {
	al := &ArtifactList{}
	path := jobIDPath(id, projectSlug) + "/artifacts"
	err := ps.client.Get(path, nil, al)
	if err != nil {
		return nil, err
	}
	return al, nil
}

// GetTestMetadata gets metadata of test in a job.
func (ps *JobOp) GetTestMetadata(id, projectSlug string) (*TestMetadataList, error) {
	tml := &TestMetadataList{}
	path := jobIDPath(id, projectSlug) + "/tests"
	err := ps.client.Get(path, nil, tml)
	if err != nil {
		return nil, err
	}
	return tml, nil
}

func jobIDPath(id, projectSlug string) string {
	return projectPathPrefix(projectSlug) + jobBasePath + "/" + id
}
